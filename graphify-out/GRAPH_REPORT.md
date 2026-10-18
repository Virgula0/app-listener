# Graph Report - app-listener  (2026-10-18)

## Corpus Check
- 466 files · ~1,644,248 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 6087 nodes · 25724 edges · 229 communities (189 shown, 40 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2392 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `3fa3325c`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- buildTrustedSet
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- rcUser
- guardUnitTest
- mkdirs
- github.com/charmbracelet/bubbles/textarea.Model
- FileTarget
- protectedPage
- github.com/charmbracelet/bubbletea.Cmd
- .launch
- Network
- IntegrationSuite
- guard_trust.bpf.c
- IntegrationSuite
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth_test.go
- globBuilder
- guard.bpf.c
- SelectOrphans
- time.Time
- icon
- RunWith
- UI
- Serve
- IntegrationSuite
- changelog.go
- picker
- github.com/spf13/cobra.Command
- IntegrationSuite
- backend
- sync.Mutex
- fileedit_model_test.go
- os.File
- IntegrationSuite
- absPath
- monitorUnitTest
- ResolvePinBase
- netGuardModel
- Open
- stdio
- editView
- fadeScroll
- LibraryClosure
- GenerateConf
- Client
- watchGroup
- runView
- TrustGuard
- runUninstall
- update_test.go
- GuardInodeKey
- Vault
- buildOneGuard
- tabs.go
- exe_supersede.h
- ListUsers
- fyne.io/fyne/v2/widget.Label
- bpfstats/main.go
- wm.go
- runMonitor
- deploy
- gui/session.go
- fcntl
- Vault
- Guard
- selection_test.go
- btrfsLayoutFrom
- integrationtests/networkguard_test.go
- tempGrantSetup
- RefreshSection
- NetEventType
- process_vm_readv.c
- GuardEvent
- DaemonEvent
- planUpdaters
- Config
- filevault_test.go
- compileGlobNames
- Finder
- KernelDev
- liveSession
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- os.FileMode
- inodeChain
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- CandidateDir
- hashCache
- openFileVaultWithMasterKey
- Well-known bypass classes
- NewDumpHook
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- PlanDaemonUnitUpdate
- Build & sign release assets (reusable workflow)
- installPage
- golang.org/x/sys/unix.Stat_t
- mountTable
- NetGuardEvent
- app-listener Terminal Demo Recording
- IntegrationSuite
- IntegrationSuite
- SampleContent
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- Monitor
- IntegrationSuite
- WebSocket /ws client
- launch_probe.c
- install.sh
- mc_tag.h
- markReady
- make test-integration (rootful Docker suite)
- .write
- supersede_probe.c
- cursor.go
- trace-app-libs.sh
- IntegrationSuite
- ResolveTempBinary
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- TempBinary
- time.Duration
- runNetworkMonitor
- Active
- catalogRefresher
- net_tester/main.go
- edit-protected
- auditGroupCoverage
- fakeGuardRepo
- daemon: permanent protection with encryption at rest
- BinaryEntry
- EventType
- layout.go
- githubRelease
- newTestTarget
- Installation
- backend_test.go
- graphify-refresh.sh
- vaultfs.go
- classifySupportError
- sealFileVaultWithMasterKey
- pinstate.go
- feed.go
- shell
- DumpHook
- WithBottomBar
- BinaryStat
- fileedit_symlink_test.go
- dialog_test.go
- Limitations
- Troubleshooting
- IntegrationSuite
- downloadFile
- Daemon configuration
- decryptStep
- secureResources
- User
- Compatibility
- How it works
- configedit_test.go
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- lockAndDeprovision
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- mountsUnitTest
- github.com/stretchr/testify/suite.Suite
- IntegrationSuite
- monitor: see what touches a path
- kernel-prof.sh
- pathCache
- .eventTypeAllowed
- Execute
- runGenKey
- writeWithin
- runNetworkGuard
- copyRegularAt
- Desktop GUI
- IntegrationSuite
- Ledger
- installSSHAgent
- changelog_test.go
- NewEventFanout
- fakeDaemon
- IntegrationSuite
- .policyTools
- Choice
- vouchedDevs
- bpfLsmListed
- bleedLayout
- Vault
- TestMain
- TestCleanRel
- VerifyInstalledResourcesLocked
- golang.org/x/sys/unix.Timespec
- MulticallError

## God Nodes (most connected - your core abstractions)
1. `Guard` - 117 edges
2. `GuardInodeKey` - 96 edges
3. `Config` - 95 edges
4. `Load()` - 89 edges
5. `IntegrationSuite` - 82 edges
6. `absPath()` - 81 edges
7. `shell` - 81 edges
8. `fileEditModel` - 77 edges
9. `writeConfig()` - 69 edges
10. `IntegrationSuite` - 68 edges

## Surprising Connections (you probably didn't know these)
- `Editor keys` --references--> `Delete()`  [INFERRED]
  docs/edit-protected.md → internal/backups/backups.go
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `swap_benign / swap_reader in-place binary swap POC` --semantically_similar_to--> `Binary ledger (/etc/app-listener/binaries.db)`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `runConfigEditor()` --calls--> `EditText()`  [INFERRED]
  cmd/functions/install/steps.go → internal/install/editor.go

## Import Cycles
- None detected.

## Communities (229 total, 40 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.09
Nodes (18): tempJournalRow, User, invalidPassword, launchRule, trustEvent, textAnswer, btrfsIoctlFsInfoArgs, infraContainerPaths() (+10 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (156): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), awaitStartupOrSignal(), loadDaemonConfig(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestAwaitStartupOrSignalAbortsAndLocksDown() (+148 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (9): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig() (+1 more)

### Community 3 - "buildTrustedSet"
Cohesion: 0.19
Nodes (13): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+5 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.08
Nodes (73): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+65 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (32): emit_event(), emit_event_kern(), fill_identity(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice() (+24 more)

### Community 7 - "rcUser"
Cohesion: 0.15
Nodes (31): rcTarget, BunTmpDir(), bunTargets(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv(), ensureRCBlock(), EnsureSSHAgentEnv(), openOwnedRegular() (+23 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (18): addErrKind, ReplacementCheck, classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest, runTool() (+10 more)

### Community 9 - "mkdirs"
Cohesion: 0.13
Nodes (36): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+28 more)

### Community 10 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.10
Nodes (39): class, pos, viewportState, classOf(), column(), dedent(), Edit(), leadingIndent() (+31 more)

### Community 11 - "FileTarget"
Cohesion: 0.06
Nodes (25): netGuardTab, policyEditor, policyHooks, policyTab, Target, newPolicyTab(), verbFor(), netWarnings() (+17 more)

### Community 12 - "protectedPage"
Cohesion: 0.07
Nodes (18): backendLink, configView, forwardForm, forwardView, protectedPage, Event, showInformation(), promptBody() (+10 more)

### Community 13 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.09
Nodes (10): diffModel, editorModel, renderBullets(), TestNoticeViewFitsAndPlacesOKLast(), TestRenderBulletsVerbatimAndIndented(), noticeModel, servedModel, serveTestModel (+2 more)

### Community 14 - ".launch"
Cohesion: 0.27
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 15 - "Network"
Cohesion: 0.17
Nodes (5): netGuardPill(), clonePolicy(), NetExtra, NetPolicy, Network

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (9): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent (+1 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (57): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+49 more)

### Community 18 - "IntegrationSuite"
Cohesion: 0.24
Nodes (6): netMonitorEvent, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (11): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, fileEditMode (+3 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.14
Nodes (49): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+41 more)

### Community 22 - "Resource"
Cohesion: 0.08
Nodes (10): ownerPaths(), Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, uniqueEncryptionRoots(), unlockRoots() (+2 more)

### Community 23 - "auth_test.go"
Cohesion: 0.13
Nodes (37): CheckPassword(), Hash(), HashConfigured(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf(), parseHash() (+29 more)

### Community 24 - "globBuilder"
Cohesion: 0.15
Nodes (12): reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName(), LibDirGlob (+4 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (72): inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), current_is_inspector() (+64 more)

### Community 26 - "SelectOrphans"
Cohesion: 0.20
Nodes (17): modifiedContextWithSource(), appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors() (+9 more)

### Community 27 - "time.Time"
Cohesion: 0.10
Nodes (17): LockedError, addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors() (+9 more)

### Community 28 - "icon"
Cohesion: 0.07
Nodes (30): binaryList, filesPage, handButton, homePage, multiAnswer, newLogo(), newBinaryList(), newBinaryRow() (+22 more)

### Community 29 - "RunWith"
Cohesion: 0.13
Nodes (18): runDiffCatalog(), needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), applyLiveRefresh(), parseMaintenanceFlags() (+10 more)

### Community 30 - "UI"
Cohesion: 0.12
Nodes (12): channelOf(), Latest(), RunWith(), TestConfirmUpdateRoutesToActiveUI(), TestRunWithRefusesConcurrentRun(), cancelUI, TestEditTextMapsUICancel(), Confirm() (+4 more)

### Community 31 - "Serve"
Cohesion: 0.08
Nodes (29): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+21 more)

### Community 33 - "changelog.go"
Cohesion: 0.08
Nodes (23): lineStyles, fg(), Highlighter, lexerFor(), NewHighlighter(), shebangLexer(), editSession(), sampleGo() (+15 more)

### Community 34 - "picker"
Cohesion: 0.10
Nodes (20): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, TestSearchTreeFollowsNoSymlink(), errText() (+12 more)

### Community 35 - "github.com/spf13/cobra.Command"
Cohesion: 0.15
Nodes (18): AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials() (+10 more)

### Community 37 - "backend"
Cohesion: 0.09
Nodes (16): checkSecret(), entryOf(), TestLiveSessionWatchAndPinger(), displayEntries(), holdForward(), pingInterval(), File, backend (+8 more)

### Community 38 - "sync.Mutex"
Cohesion: 0.08
Nodes (25): readConfigPut(), TestReadGrantRequestConfig(), newControlManager(), readAuthRequest(), readBareRequest(), readClientLines(), readForwardRequest(), readGrantRequest() (+17 more)

### Community 39 - "fileedit_model_test.go"
Cohesion: 0.21
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 40 - "os.File"
Cohesion: 0.04
Nodes (70): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir() (+62 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.07
Nodes (15): IntegrationSuite, decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest (+7 more)

### Community 42 - "absPath"
Cohesion: 0.07
Nodes (6): shQuote(), IntegrationSuite, inNosuidNamespace(), IntegrationSuite, IntegrationSuite, absPath()

### Community 44 - "ResolvePinBase"
Cohesion: 0.50
Nodes (5): dirIsRootOwnedSafe(), isBpffs(), mountBpffs(), ResolvePinBase(), TestDirIsRootOwnedSafe()

### Community 45 - "netGuardModel"
Cohesion: 0.05
Nodes (37): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+29 more)

### Community 46 - "Open"
Cohesion: 0.18
Nodes (19): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+11 more)

### Community 47 - "stdio"
Cohesion: 0.07
Nodes (4): put(), wait_then_read(), hold_forever(), main()

### Community 48 - "editView"
Cohesion: 0.12
Nodes (8): codeEntry, editView, Entry, selectable(), fileSize(), joinRel(), newCodeEntry(), relOf()

### Community 49 - "fadeScroll"
Cohesion: 0.07
Nodes (16): fadeBar, fadeBarRenderer, fadeScroll, fadeScrollRenderer, fadeScrollOf(), newFadeScroll(), newTestFadeScroll(), settle() (+8 more)

### Community 50 - "LibraryClosure"
Cohesion: 0.06
Nodes (35): sanityCheckBinary(), TestSanityCheckBinaryRefusesOtherMachine(), ldConfParser, Check(), buildTiny(), TestCheck(), defaultLibDirs(), elfInterp() (+27 more)

### Community 51 - "GenerateConf"
Cohesion: 0.05
Nodes (55): TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), discordCatalogEntry(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles(), TestGeneratedLibraryBlockRoundTrips() (+47 more)

### Community 52 - "Client"
Cohesion: 0.08
Nodes (16): reserveStdout(), runBackend(), Client, launch(), TestHelperBackend(), TestReaderCapsLines(), Hello, Reader (+8 more)

### Community 53 - "watchGroup"
Cohesion: 0.11
Nodes (35): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+27 more)

### Community 54 - "runView"
Cohesion: 0.08
Nodes (10): logGrid, logLine, runView, stageItem, stageState, asset(), coloredIcon(), logRow() (+2 more)

### Community 55 - "TrustGuard"
Cohesion: 0.06
Nodes (17): tempGrant, tempMaskOp, tempRow, trustHook, trustRows, deleteResKeys(), setSupersedeTrust(), liveTrust() (+9 more)

### Community 56 - "runUninstall"
Cohesion: 0.18
Nodes (19): deletePostBackups(), listBackups(), restoreBackups(), offerBackupCleanup(), runUninstall(), removeKeyAndEmptyDir(), removeMasterKey(), TestRemoveKeyAndEmptyDir() (+11 more)

### Community 57 - "update_test.go"
Cohesion: 0.12
Nodes (11): fetchReleases(), InstalledVersion(), readInstalledVersion(), TestFetchReleases(), TestFetchReleasesHTTPError(), TestFilterChannel(), TestNewerThanStable(), TestParseStableVersion() (+3 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.03
Nodes (50): engine, refusedReplacements, resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), TestSetInspectorsBounded(), appendInodeRow() (+42 more)

### Community 59 - "Vault"
Cohesion: 0.17
Nodes (14): isFileVaultCiphertext(), isRegularFileTarget(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr(), isNotEncryptedErr(), newBoundedKeyFn(), readKey() (+6 more)

### Community 60 - "buildOneGuard"
Cohesion: 0.05
Nodes (65): CheckBPFLSM(), CheckEBPF(), backingDeviceUnion(), buildConcurrency(), buildGuards(), buildOneGuard(), eventFilterOptions(), newPinGeneration() (+57 more)

### Community 61 - "tabs.go"
Cohesion: 0.06
Nodes (9): navItem, navItemRenderer, tabBar, tabButton, tabButtonRenderer, tabItem, newNavItem(), newTabBar() (+1 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.14
Nodes (20): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+12 more)

### Community 63 - "ListUsers"
Cohesion: 0.05
Nodes (42): ParseEventsFlag(), auditAfterEdit(), dialLiveSession(), LiveModeAvailable(), TestValidateFlagsEditConfig(), applyEdit(), checkAndApply(), editAgain() (+34 more)

### Community 64 - "fyne.io/fyne/v2/widget.Label"
Cohesion: 0.05
Nodes (16): appTheme, binaryRow, clickPanel, iconTile, iconTileRenderer, panel, panelRenderer, pill (+8 more)

### Community 65 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 66 - "wm.go"
Cohesion: 0.11
Nodes (25): detectWM(), reachable(), centered(), Detect(), FitSize(), Session, hyprlandSocket(), i3RuntimeSocket() (+17 more)

### Community 67 - "runMonitor"
Cohesion: 0.15
Nodes (15): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+7 more)

### Community 68 - "deploy"
Cohesion: 0.24
Nodes (14): applyUpdate(), assetsFor(), checkUpdatePreconditions(), deploy(), downloadAndVerify(), FetchLatest(), FetchVerified(), resolveAssets() (+6 more)

### Community 69 - "gui/session.go"
Cohesion: 0.20
Nodes (15): ancestorEnv(), invokingHome(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField(), TestDisplayValue() (+7 more)

### Community 70 - "fcntl"
Cohesion: 0.09
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 71 - "Vault"
Cohesion: 0.24
Nodes (3): Vault, Open(), TestOpenVaultRefusesSymlink()

### Community 72 - "Guard"
Cohesion: 0.05
Nodes (24): lockVaultFully(), vaultOpUnderGuard(), deferredBinary, exeRow, rootHandle, openRoot(), SharedPinDegraded(), TestNoteRootMountFollowsMountID() (+16 more)

### Community 73 - "selection_test.go"
Cohesion: 0.23
Nodes (10): clipboardText(), setClipboard(), cellOf(), repeat(), selKey(), TestSelectionEditing(), TestSelectionLargeDelete(), TestSelectionMouseDrag() (+2 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 75 - "integrationtests/networkguard_test.go"
Cohesion: 0.27
Nodes (5): IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 76 - "tempGrantSetup"
Cohesion: 0.29
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 77 - "RefreshSection"
Cohesion: 0.09
Nodes (30): logRefreshChange(), sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange() (+22 more)

### Community 78 - "NetEventType"
Cohesion: 0.06
Nodes (20): snapNetMonitor, NetBpfEvent, NetEvent, NetEventType, eventsOrAll(), newNetEventChecks(), FormatAddr(), NetEventTypes() (+12 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.18
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "GuardEvent"
Cohesion: 0.06
Nodes (23): holdHeadless(), HeadlessLine(), runGuardHeadless(), bpfGuardEvent, snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestDispatchWithoutOwnAllowedEvents() (+15 more)

### Community 81 - "DaemonEvent"
Cohesion: 0.07
Nodes (42): ServeConfig, UIDResolver, NewUIDResolver(), catchLifecycleSignals(), drainEphemeral(), newDaemonModel(), runDaemon(), runDaemonUI() (+34 more)

### Community 82 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 83 - "Config"
Cohesion: 0.09
Nodes (39): newCatalogRefresher(), makeReloadHandler(), noteReloadedSettings(), startCatalogRefresh(), resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted(), startTrustGuard(), trustStartupError() (+31 more)

### Community 84 - "filevault_test.go"
Cohesion: 0.24
Nodes (20): looksLikeFileVaultRecord(), inode(), TestIsEncryptedDispatchesToFileVaultForRegularFiles(), TestIsFileVaultRecordShape(), TestIsProvisionedForFile(), TestLockUnlockFileInPlaceIdempotent(), TestUnlockFileInPlaceRecoversFromInterruptedTransform(), TestUnlockFileInPlaceWrongKeyFails() (+12 more)

### Community 85 - "compileGlobNames"
Cohesion: 0.17
Nodes (9): GlobChild, statFunc, compileGlobNames(), globText(), resolveBits(), resolveChildren(), statKey(), TestCompileGlobNames_RepeatedPatternOrsBits() (+1 more)

### Community 86 - "Finder"
Cohesion: 0.15
Nodes (19): findButton, span, findAll(), Finder, lowerRunes(), NewFinder(), newFindInput(), overlay() (+11 more)

### Community 87 - "KernelDev"
Cohesion: 0.15
Nodes (17): fdStat, physicalParent(), statFD(), KernelDev(), BtrfsDevices(), btrfsSysfsDevices(), TestBtrfsSysfsDevices(), codeVouched() (+9 more)

### Community 88 - "liveSession"
Cohesion: 0.26
Nodes (8): parseErr(), TestLiveSessionConfigAndPut(), chooseResources(), eventsLabel(), forwardBinaries(), normalizedEvents(), runForward(), liveSession

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "os.FileMode"
Cohesion: 0.15
Nodes (18): installConfig(), installFile(), installFileAs(), installUnit(), readAtNoFollow(), upsertFile(), upsertUnitAt(), writeAtNoFollow() (+10 more)

### Community 92 - "inodeChain"
Cohesion: 0.26
Nodes (13): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+5 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.09
Nodes (11): foreignPinnedLink(), eventTypeKey(), mcNameKey(), statInodeKey(), judgePasses(), parseProcs(), TestJudgePasses(), TestParseProcs() (+3 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "CandidateDir"
Cohesion: 0.08
Nodes (27): cleanOrphanedFscrypt(), confSafeMatch(), Discover(), DiscoverForUsers(), DiscoverSystem(), expandPlaceholders(), CandidateDir, homeMatchConfined() (+19 more)

### Community 100 - "hashCache"
Cohesion: 0.19
Nodes (7): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache

### Community 101 - "openFileVaultWithMasterKey"
Cohesion: 0.26
Nodes (13): classifyRegularFileTarget(), clearRecovery(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), openFileVaultWithMasterKey(), recoverFileInPlace(), rejectSidecarSymlink(), stageRecovery() (+5 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "NewDumpHook"
Cohesion: 0.24
Nodes (11): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+3 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "PlanDaemonUnitUpdate"
Cohesion: 0.19
Nodes (13): confirmUnit(), TestUnifiedDiff(), TestUnifiedDiffIdentical(), UnifiedDiff(), ConfirmUnitFunc, InstalledDaemonUnit(), PlanDaemonUnitUpdate(), ShippedDaemonUnit() (+5 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "installPage"
Cohesion: 0.08
Nodes (18): binSource, installAction, installedBinary, InstallMode, installPage, InstallRequest, latestRelease, Options (+10 more)

### Community 110 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 111 - "mountTable"
Cohesion: 0.27
Nodes (7): MountInfo, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), mountTable

### Community 112 - "NetGuardEvent"
Cohesion: 0.11
Nodes (16): snapNetGuard, NetGuardEvent, NetworkGuardRepository, NewGuardUseCase(), NewMonitorUseCase(), NetworkGuardUseCase, NewNetworkGuardUseCase(), newFakeNetworkGuardRepo() (+8 more)

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 116 - "SampleContent"
Cohesion: 0.15
Nodes (15): TestDaemonUnitWithMetadataOutput(), detectSSHAgentUnits(), isInstallerSSHAgentUnit(), removeSSHAgentEnv(), removeSSHAgentUnit(), revertSSHAgents(), staleRCUsers(), TestIsInstallerSSHAgentUnit() (+7 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "Monitor"
Cohesion: 0.15
Nodes (4): NewMonitor(), newMountTable(), Monitor, watchedInode

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "launch_probe.c"
Cohesion: 0.53
Nodes (4): copy(), is_verb(), main(), usage()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.27
Nodes (9): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netm_exec_applet() (+1 more)

### Community 127 - "markReady"
Cohesion: 0.25
Nodes (9): markReady(), notifySystemdReady(), inodeBatchFn, fakeBatches(), TestBatchInodeRows(), countInodeRows(), InodeOccupancy(), TestCountInodeRows() (+1 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - ".write"
Cohesion: 0.11
Nodes (17): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), testKeyPair(), TestParseChecksum(), TestParsePublicKey(), TestVerifyRelease() (+9 more)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "cursor.go"
Cohesion: 0.14
Nodes (13): buttonCursor, handCursors(), handCursorsOn(), newButtonCursor(), showConfirm(), showCustomConfirm(), showDialog(), showErrorDialog() (+5 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "ResolveTempBinary"
Cohesion: 0.18
Nodes (8): confirmUserPlaced(), forwardCheck(), binaryStatOf(), loadPinnedExeMaps(), ResolveTempBinary(), stripPinnedAllow(), StripPinnedTempAllows(), TestResolveTempBinaryClassesRuntime()

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "TempBinary"
Cohesion: 0.12
Nodes (16): TempBinary, TemporaryGrant, TempRule, planTempBlock(), applyTempGrants(), closeTempBinaries(), TemporaryAccess, logTempGrant() (+8 more)

### Community 143 - "time.Duration"
Cohesion: 0.14
Nodes (13): procStats, statsSample, cpuTime(), humanBytes(), rssBytes(), New(), recv(), TestSubDeliversMissedTickLate() (+5 more)

### Community 144 - "runNetworkMonitor"
Cohesion: 0.13
Nodes (10): ParseNetEventsFlag(), runHeadless(), runHeadless(), runNetworkMonitor(), runTUI(), ProtocolString(), NetworkMonitorRepository, NewNetModel() (+2 more)

### Community 145 - "Active"
Cohesion: 0.21
Nodes (10): nextConfig(), showChangelog(), confirmUpdate(), isTerminal(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), EditText(), Active(), Field (+2 more)

### Community 146 - "catalogRefresher"
Cohesion: 0.15
Nodes (7): drain(), readEvents(), swapFile(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), catalogRefresher, debouncer

### Community 147 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 148 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 149 - "auditGroupCoverage"
Cohesion: 0.20
Nodes (12): auditEntry(), auditFindings(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths() (+4 more)

### Community 152 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 153 - "BinaryEntry"
Cohesion: 0.08
Nodes (40): binaryVerifyState, HashStore, sharedHashEntry, BinaryEntry, checkBtrfsKeys(), BinariesSummary(), canonicalBinaryPath(), canonicalPaths() (+32 more)

### Community 154 - "EventType"
Cohesion: 0.08
Nodes (15): snapEngines, snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting() (+7 more)

### Community 155 - "layout.go"
Cohesion: 0.23
Nodes (10): adaptiveGrid, checkGrid, gapVBox, TestAdaptiveGridDropsColumns(), cellMin(), newAdaptiveGrid(), newCheckGrid(), newGapVBox() (+2 more)

### Community 156 - "githubRelease"
Cohesion: 0.19
Nodes (21): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), tags(), TestReleasesSince(), compareStableVersions(), filterChannel() (+13 more)

### Community 157 - "newTestTarget"
Cohesion: 0.11
Nodes (21): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+13 more)

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "backend_test.go"
Cohesion: 0.28
Nodes (18): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+10 more)

### Community 161 - "vaultfs.go"
Cohesion: 0.14
Nodes (21): DefaultNewFileMode(), FileModeToUnixMode(), IsBinaryContent(), IsOctalMode(), inodeOf(), mustRead(), mustWrite(), openRootT() (+13 more)

### Community 162 - "classifySupportError"
Cohesion: 0.17
Nodes (10): classifySetupError(), classifySupportError(), TestClassifySetupErrorGeneric(), TestClassifySetupErrorNotSetup(), TestClassifySetupErrorNotSupported(), TestClassifySupportErrorEncryptionNotEnabled(), TestClassifySupportErrorEncryptionNotEnabledF2fs(), TestClassifySupportErrorGeneric() (+2 more)

### Community 163 - "sealFileVaultWithMasterKey"
Cohesion: 0.24
Nodes (11): deriveFileVaultSubkey(), newFileVaultAEAD(), openFileVault(), sealFileVault(), sealFileVaultWithMasterKey(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestOpenFileVaultRejectsTampering() (+3 more)

### Community 164 - "pinstate.go"
Cohesion: 0.13
Nodes (28): configureDaemonLogging(), lockOneRoot(), lockRootRecovering(), relockStaleVaults(), resolveConfigPath(), runLockdown(), ensurePinStateFilePlaceholder(), ensurePlaceholder() (+20 more)

### Community 165 - "feed.go"
Cohesion: 0.06
Nodes (45): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+37 more)

### Community 166 - "shell"
Cohesion: 0.05
Nodes (37): installUI, pageDef, pageID, resizeWatch, shell, sizedDialog, sysStatus, newShell() (+29 more)

### Community 167 - "DumpHook"
Cohesion: 0.24
Nodes (5): logHook, formatEntry(), sortedKeys(), sortStrings(), DumpHook

### Community 168 - "WithBottomBar"
Cohesion: 0.31
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 169 - "BinaryStat"
Cohesion: 0.29
Nodes (10): BinaryStat, binaryStatOf(), StatBinary(), StatBinaryFile(), durableFS(), fsidOf(), hostVouches(), StatDurable() (+2 more)

### Community 170 - "fileedit_symlink_test.go"
Cohesion: 0.36
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 171 - "dialog_test.go"
Cohesion: 0.45
Nodes (10): openPicker(), pickerTree(), shownNames(), TestPickerArrowHighlightsFolderEnterOpensIt(), TestPickerCtrlFFocusesSearch(), TestPickerFolderModeChoosesShownFolder(), TestPickerSearchFiltersFolder(), TestPickerSubfolderSearchPicksFile() (+2 more)

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 175 - "downloadFile"
Cohesion: 0.24
Nodes (7): downloadFile(), downloadReleaseFiles(), fetchRelease(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleaseByTag(), progressReader

### Community 176 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Several guarded trees, one vault, Watch sections

### Community 177 - "decryptStep"
Cohesion: 0.22
Nodes (9): askDecrypt(), cleanOrphanedMetadata(), decryptDirectories(), decryptStep(), pickDirsToDecrypt(), TestDecryptDirectoriesEmpty(), CleanOrphanedMetadata(), LivePolicyDescriptors() (+1 more)

### Community 178 - "secureResources"
Cohesion: 0.21
Nodes (9): encryptDirectories(), secureResources(), verifyEncryptionState(), confirmRunPrereq(), resolveFilesystemPrereqs(), askEncryption(), harnessSuite, EnsureSystemSetup() (+1 more)

### Community 179 - "User"
Cohesion: 0.05
Nodes (58): entryWriters(), askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir(), applyDiffAdditions(), pathCovered() (+50 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 182 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 186 - "Binary trust"
Cohesion: 0.29
Nodes (7): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 189 - "lockAndDeprovision"
Cohesion: 0.12
Nodes (12): deprovisionKind, migrationTarget, writeTempWithMetadata(), classifyDeprovision(), classifyDeprovisionErr(), classifyMigrationTarget(), isDeprovisionBusy(), isDeprovisionMissing() (+4 more)

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "app-listener documentation"
Cohesion: 0.40
Nodes (5): app-listener documentation, Getting started, Modes, Operating, The daemon in depth

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.50
Nodes (4): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network-monitor: see what a program talks to

### Community 195 - "github.com/stretchr/testify/suite.Suite"
Cohesion: 0.20
Nodes (3): mountinfoUnitTest, unescapeMountinfo(), editorHarness

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 200 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 201 - "Execute"
Cohesion: 0.18
Nodes (10): installHooks(), releaseBuild(), Requested(), Run(), TestRequested(), TestSubcommandOutranksGUI(), Execute(), initLogger() (+2 more)

### Community 202 - "runGenKey"
Cohesion: 0.20
Nodes (11): confirmMasterKeyOverwrite(), runGenKey(), runOneShotMode(), ensureInstalledBinary(), ensureMasterKey(), prepareInstallation(), TestEnsureInstalledBinary(), checkKeyLen() (+3 more)

### Community 203 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 204 - "runNetworkGuard"
Cohesion: 0.09
Nodes (26): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), snapNet, fillNetwork() (+18 more)

### Community 205 - "copyRegularAt"
Cohesion: 0.31
Nodes (10): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), fchownFromInfo(), futimesFromInfo(), readDirNames() (+2 more)

### Community 206 - "Desktop GUI"
Cohesion: 0.17
Nodes (12): Actions, Desktop GUI, File and folder pickers, Files, Home, Install, Network, Protected (+4 more)

### Community 209 - "Ledger"
Cohesion: 0.10
Nodes (12): candidates(), confirmable(), confirmAll(), describe(), hashNow(), run(), fyneLogFilter, applySchema() (+4 more)

### Community 211 - "installSSHAgent"
Cohesion: 0.29
Nodes (7): dupFD(), installSSHAgent(), installSSHAgentEnv(), TestInstallSSHAgentRefusesSymlinkedUnitDir(), EnsureBunTmpDir(), TestEnsureBunTmpDirCreates0700(), DescendCreateOwned()

### Community 212 - "changelog_test.go"
Cohesion: 0.31
Nodes (8): newChangelogModel(), rel(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog(), TestNewChangelogModelViewContainsTitleAndNotes()

### Community 213 - "NewEventFanout"
Cohesion: 0.24
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 217 - "Choice"
Cohesion: 0.50
Nodes (3): MultiSelectKeymap(), Choice, MultiSelect()

### Community 218 - "vouchedDevs"
Cohesion: 0.33
Nodes (6): isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed(), vouchedDevs()

### Community 219 - "bpfLsmListed"
Cohesion: 0.40
Nodes (6): bpfLsmListed(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 222 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

### Community 224 - "VerifyInstalledResourcesLocked"
Cohesion: 0.67
Nodes (3): TestVerifyResourcesLocked(), VerifyInstalledResourcesLocked(), VerifyResourcesLocked()

## Knowledge Gaps
- **32 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC`, `statonly stat/statx metadata leak`, `app-listener-gui best-effort build` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 565 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **40 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Editor keys` connect `edit-protected` to `runUninstall`, `backend`?**
  _High betweenness centrality (0.027) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.08712125243452071 - nodes in this community are weakly interconnected._
- **Why does `edit-protected` connect `edit-protected` to `README.md`?**
  _High betweenness centrality (0.027) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.01721698113207547 - nodes in this community are weakly interconnected._
- **Why does `absPath()` connect `absPath` to `daemon/daemon.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `integrationtests/networkguard_test.go`, `.launch`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `TestMain`?**
  _High betweenness centrality (0.026) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.08957708957708958 - nodes in this community are weakly interconnected._