# Graph Report - app-listener  (2026-10-30)

## Corpus Check
- 562 files · ~1,789,339 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 7315 nodes · 33414 edges · 253 communities (208 shown, 45 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2995 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `a6e2b8ca`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- go_pkg_strings
- testing.T
- IntegrationSuite
- op
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- IntegrationSuite
- guardUnitTest
- mkdirs
- GuardEvent
- safeio.go
- editView
- prepareDaemonStart
- Network
- Candidate
- IntegrationSuite
- guard_trust.bpf.c
- fyne.io/fyne/v2/widget.Label
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- FileTarget
- guard.bpf.c
- LibraryClosure
- RunCmd
- updatePage
- desktop.go
- time.Time
- TestMonitorUseCaseLifecycle
- fileedit_model_test.go
- highlight_test.go
- picker
- secret_probe/main.go
- IntegrationSuite
- Message
- controlServer
- render.go
- .refreshOnce
- IntegrationSuite
- absPath
- monitorUnitTest
- cursor.go
- binaryVetter
- Guard
- stdio
- update_test.go
- fadeScroll
- CandidateDir
- SanitizeText
- Client
- Config
- session
- .apply
- changelog.go
- mutedText
- engine
- Monitor
- Prereq
- .FilterExistingWhitelist
- io.Writer
- .filterEvent
- policyEditor
- unix_connect.c
- wm.go
- IntegrationSuite
- linkView
- NetGuardEvent
- Vault
- daemon_multicall_test.go
- journal
- fcntl
- BinaryEntry
- buildTrustedSet
- codecheck/main.go
- runMonitor
- listFacultatives
- process_vm_readv.c
- GuardInodeKey
- recordFlow
- sync.Mutex
- rcUser
- OpenSystemPlaced
- bpfstats/main.go
- exe_supersede.h
- loginServer
- vaultfs.go
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- DumpHook
- tempGrantSetup
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- gui/session.go
- hashCache
- Well-known bypass classes
- writeWithin
- jit_provenance.c
- ptrace_seize.c
- recordops.c
- Open
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- bufio.Reader
- ConfinedEntry
- UI
- KernelDev
- IntegrationSuite
- maintainerKeeper
- PlanDaemonUnitUpdate
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- .netIntegrityContainer
- Gate
- WebSocket /ws client
- inet_connect.c
- install.sh
- mc_tag.h
- classify
- make test-integration (rootful Docker suite)
- newTestTarget
- catalogRefresher
- planUpdaters
- trace-app-libs.sh
- IntegrationSuite
- Finder
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- Daemon configuration
- serve.go
- liveSession
- github.com/charmbracelet/bubbles/textarea.Model
- Ledger
- daemon: permanent protection with encryption at rest
- fakeGuardRepo
- os.FileInfo
- find
- fakeStore
- EventType
- IntegrationSuite
- buildOneGuard
- watchSet
- copyRegularAt
- sources
- Installation
- backend_test.go
- graphify-refresh.sh
- .Layout
- NetEventType
- .populateWatchInodes
- promptPassword
- feed.go
- shell
- sources_test.go
- RunWith
- Run
- SelectOrphans
- Limitations
- runOfflineEdit
- IntegrationSuite
- TrustGuard
- net_tester/main.go
- IntegrationSuite
- update/update.go
- selection_test.go
- IntegrationSuite
- How it works
- inodeChain
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- github.com/spf13/cobra.Command
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- IntegrationSuite
- openBunTmp
- edit-protected
- inspector_probe.c
- kernel-prof.sh
- filevault_test.go
- pathCache
- pinstate.go
- DaemonEvent
- WithBottomBar
- runBackend
- batchInodeRows
- monitor
- IntegrationSuite
- revertSystemFiles
- main_test.go
- .eventTypeAllowed
- BusClient
- procstats.go
- writeEvent
- .awaitDaemonUpOrDead
- .TestRecord_GuardParity
- SecretStoreNotes
- parseDaemonEvents
- saslServer
- WithSuspectsFrom
- maintainer/plan.go
- ValidateTargets
- ReadInstallCert
- downloadFile
- supersede_probe.c
- github.com/charmbracelet/bubbletea.Cmd
- AddServeFlags
- mountTable
- launchKey
- github.com/Virgula0/app-listener/tools/codecheck
- fileedit_symlink_test.go
- TestReadFrame
- DiscoverInfraBinaries
- bleedLayout
- inode_dev.h
- Troubleshooting
- VouchedDevs
- GenerateConf
- parseReservedPorts
- shQuote
- taint_victim.c
- User
- parseKernelRelease
- Check
- golang.org/x/sys/unix.Timespec
- fileID
- monitor: see what touches a path
- MulticallError
- plan_test.go

## God Nodes (most connected - your core abstractions)
1. `Config` - 146 edges
2. `Guard` - 125 edges
3. `GuardInodeKey` - 112 edges
4. `shell` - 97 edges
5. `Load()` - 95 edges
6. `absPath()` - 93 edges
7. `IntegrationSuite` - 83 edges
8. `fileEditModel` - 77 edges
9. `writeConfig()` - 74 edges
10. `IntegrationSuite` - 67 edges

## Surprising Connections (you probably didn't know these)
- `eBPF at three levels` --references--> `monitor()`  [INFERRED]
  docs/how-it-works.md → integrationtests/secret_probe/main.go
- `Why not an AppArmor profile?` --references--> `monitor()`  [INFERRED]
  docs/how-it-works.md → integrationtests/secret_probe/main.go
- `Distributions` --references--> `monitor()`  [INFERRED]
  docs/compatibility.md → integrationtests/secret_probe/main.go
- `Files` --references--> `monitor()`  [INFERRED]
  docs/gui.md → integrationtests/secret_probe/main.go
- `Workspaces` --references--> `monitor()`  [INFERRED]
  docs/gui.md → integrationtests/secret_probe/main.go

## Import Cycles
- None detected.

## Communities (253 total, 45 thin omitted)

### Community 0 - "go_pkg_strings"
Cohesion: 0.11
Nodes (14): tempJournalRow, invalidPassword, recordTick, linkedFile, argvMatch, launchRule, trustEvent, btrfsIoctlFsInfoArgs (+6 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (168): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+160 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (7): configBinaryLines(), IntegrationSuite, probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig()

### Community 3 - "op"
Cohesion: 0.17
Nodes (13): clearExcept(), createTemp(), holds(), readlinkat(), removeDir(), ignoreNotExist(), openBeneath(), openDir() (+5 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (78): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+70 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.09
Nodes (51): prctl_gate(), guard_task_prctl(), admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader() (+43 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.11
Nodes (48): current_exe(), emit_dentry(), emit_event(), emit_event_kern(), emit_file_write(), emit_kern(), emit_unix_peer(), emit_user() (+40 more)

### Community 7 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 8 - "guardUnitTest"
Cohesion: 0.05
Nodes (26): addErrKind, ReplacementCheck, classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest, runTool() (+18 more)

### Community 9 - "mkdirs"
Cohesion: 0.09
Nodes (41): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+33 more)

### Community 10 - "GuardEvent"
Cohesion: 0.06
Nodes (25): displayEntries(), holdHeadless(), HeadlessLine(), resolveGuardConfig(), runGuard(), runGuardHeadless(), runGuardTUI(), bpfGuardEvent (+17 more)

### Community 11 - "safeio.go"
Cohesion: 0.12
Nodes (26): dupFD(), installFileAs(), installSSHAgent(), installSSHAgentEnv(), installUnit(), readAtNoFollow(), TestInstallSSHAgentRefusesSymlinkedUnitDir(), upsertFile() (+18 more)

### Community 12 - "editView"
Cohesion: 0.11
Nodes (10): codeEntry, editView, Entry, selectable(), fileSize(), joinRel(), newCodeEntry(), newEditView() (+2 more)

### Community 13 - "prepareDaemonStart"
Cohesion: 0.06
Nodes (37): CheckBPFLSM(), CheckEBPF(), configureDaemonLogging(), confirmMasterKeyOverwrite(), loadDaemonConfig(), newPinGeneration(), noteReloadedSettings(), prepareDaemonStart() (+29 more)

### Community 14 - "Network"
Cohesion: 0.14
Nodes (8): netGuardPill(), NetworkGuardRepository, clonePolicy(), NetEngines, NetExtra, NetPolicy, Network, NetGuardConfig

### Community 15 - "Candidate"
Cohesion: 0.07
Nodes (39): appendSectionsAndEdit(), collectDiffAdditions(), TestDiffMergePreservesExistingAndParses(), TestUncoveredCandidates(), uncoveredCandidates(), selectAndEditConfig(), addManualDirectories(), editConfig() (+31 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (8): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (64): argv_env(), argv_matched(), argv_rule_at(), argv_word(), caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key() (+56 more)

### Community 18 - "fyne.io/fyne/v2/widget.Label"
Cohesion: 0.03
Nodes (24): appTheme, binaryRow, clickPanel, iconTile, iconTileRenderer, navItem, navItemRenderer, panel (+16 more)

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (12): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, dirEntry (+4 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.15
Nodes (48): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+40 more)

### Community 22 - "Resource"
Cohesion: 0.03
Nodes (39): ownerPaths(), newControlManager(), startControlManager(), backingDeviceUnion(), buildGuards(), makeReloadHandler(), reloadOnce(), reloadSlotsNeeded() (+31 more)

### Community 23 - "auth.go"
Cohesion: 0.12
Nodes (41): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+33 more)

### Community 24 - "FileTarget"
Cohesion: 0.14
Nodes (9): Files, FileSpec, FileTarget, Policy, newFileTarget(), NormalizeSpec(), overlaps(), stopPumped() (+1 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.07
Nodes (93): inode_dev(), add_inode_to_guard(), admit_port(), btrfs_copy_gate(), bus_key(), check_and_emit_args(), chmod_only_drops_write(), code_suspect() (+85 more)

### Community 26 - "LibraryClosure"
Cohesion: 0.07
Nodes (29): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), IsScript(), ldConfInclude(), ldSoConfDirs() (+21 more)

### Community 27 - "RunCmd"
Cohesion: 0.07
Nodes (45): offerDaemonStart(), runDiffCatalog(), applyLiveRefresh(), buildBinaryIfNeeded(), certifyInstall(), installBinaryOnly(), markInstallComplete(), preflightDeployedBinary() (+37 more)

### Community 28 - "updatePage"
Cohesion: 0.14
Nodes (15): installedBinary, InstallMode, InstallRequest, latestRelease, Options, UninstallHooks, UninstallRequest, UpdateHooks (+7 more)

### Community 29 - "desktop.go"
Cohesion: 0.04
Nodes (48): releaseBuild(), Run(), TestEveryInstallModeMapped(), Progress, User, fakeLink, logGrid, logHook (+40 more)

### Community 30 - "time.Time"
Cohesion: 0.05
Nodes (50): LockedError, lines(), pickBinaries(), pickLibs(), reviewRecording(), sweepID, sweepMemo, appletLine() (+42 more)

### Community 31 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.18
Nodes (12): NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestMonitorUseCaseStartError() (+4 more)

### Community 32 - "fileedit_model_test.go"
Cohesion: 0.20
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 33 - "highlight_test.go"
Cohesion: 0.13
Nodes (27): New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep(), TestCursorDoesNotBlink(), TestEditKeys(), TestEnterKeepsIndentation(), TestEnterPastNinetyNineLines(), TestSetTextStartsAtFirstLine() (+19 more)

### Community 34 - "picker"
Cohesion: 0.08
Nodes (31): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, lstatOf(), outsideHardlinks() (+23 more)

### Community 35 - "secret_probe/main.go"
Cohesion: 0.26
Nodes (23): busPath(), connexec(), connFD(), createItem(), dial(), execInto(), failf(), get() (+15 more)

### Community 37 - "Message"
Cohesion: 0.09
Nodes (14): checkSecret(), checkTrustLine(), trustRows(), File, TrustItem, backend, backendOp, backendState (+6 more)

### Community 38 - "controlServer"
Cohesion: 0.17
Nodes (9): authSelfPeer(), readClientLines(), startControlServer(), streamEvents(), TestEditControlSessionEndIdempotent(), TestRecordFailureLockout(), configEditor, controlServer (+1 more)

### Community 39 - "render.go"
Cohesion: 0.12
Nodes (13): lineStyles, viewportState, baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces(), setGutter() (+5 more)

### Community 40 - ".refreshOnce"
Cohesion: 0.53
Nodes (5): discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad()

### Community 41 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 42 - "absPath"
Cohesion: 0.07
Nodes (5): IntegrationSuite, IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 44 - "cursor.go"
Cohesion: 0.09
Nodes (19): buttonCursor, handCheckGroup, handRadioGroup, handSelect, handCursor(), handCursors(), handCursorsOn(), newButtonCursor() (+11 more)

### Community 45 - "binaryVetter"
Cohesion: 0.08
Nodes (22): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), systemPatterns(), systemRule(), resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted(), startTrustGuard() (+14 more)

### Community 46 - "Guard"
Cohesion: 0.04
Nodes (31): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), deferredBinary, exeRow, rootHandle, supersededBinary (+23 more)

### Community 48 - "update_test.go"
Cohesion: 0.08
Nodes (23): fetchReleases(), newerThanInstalled(), parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFetchReleases(), TestFetchReleasesHTTPError() (+15 more)

### Community 49 - "fadeScroll"
Cohesion: 0.08
Nodes (14): fadeBar, fadeBarRenderer, fadeScroll, fadeScrollRenderer, newTestFadeScroll(), settle(), TestFadeBarDragEndsGlide(), TestFadeBarDragMovesContent() (+6 more)

### Community 50 - "CandidateDir"
Cohesion: 0.07
Nodes (27): entryWriters(), reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName() (+19 more)

### Community 51 - "SanitizeText"
Cohesion: 0.05
Nodes (27): backendLink, configView, forwardForm, forwardView, installUI, protectedPage, recordReview, recordView (+19 more)

### Community 52 - "Client"
Cohesion: 0.10
Nodes (12): Client, childEnv(), launch(), Start(), verifiedBinary(), TestChildEnvIsAllowlisted(), TestHelperBackend(), TestReaderCapsLines() (+4 more)

### Community 53 - "Config"
Cohesion: 0.06
Nodes (75): sameConfig(), ensureAnchor(), ensurePortAnchors(), ensureRootDir(), ephemeralRange(), ensureSecretAnchors(), auditGroupCoverage(), findResource() (+67 more)

### Community 54 - "session"
Cohesion: 0.12
Nodes (9): peerCred(), TestClientAuth(), TestSASL(), clientAuth(), readFrame(), serverAuth(), queued, senderQueue (+1 more)

### Community 55 - ".apply"
Cohesion: 0.14
Nodes (10): TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(), TestApplyRollsBackOnRefusal(), testPlan(), applier, ledgerOps, PlannedBinary (+2 more)

### Community 56 - "changelog.go"
Cohesion: 0.17
Nodes (22): changelogText(), cumulativeChangelog(), newChangelogModel(), newReleaseWindow(), releasesSince(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit() (+14 more)

### Community 57 - "mutedText"
Cohesion: 0.05
Nodes (49): actionGate, binaryList, binSource, checkGrid, filesPage, handButton, handCheck, installAction (+41 more)

### Community 58 - "engine"
Cohesion: 0.04
Nodes (31): engine, portOwner, resInfo, taintOwner, TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), TestSetInspectorsBounded() (+23 more)

### Community 59 - "Monitor"
Cohesion: 0.18
Nodes (4): NewMonitor(), putBtrfsLayout(), Monitor, watchedInode

### Community 60 - "Prereq"
Cohesion: 0.20
Nodes (6): classifySetupError(), TestClassifySetupErrorGeneric(), TestClassifySetupErrorNotSetup(), TestClassifySetupErrorNotSupported(), Prereq, fakeVault

### Community 61 - ".FilterExistingWhitelist"
Cohesion: 0.10
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 62 - "io.Writer"
Cohesion: 0.13
Nodes (16): formatElapsed(), orNone(), printBinary(), printRecording(), candidates(), confirmAll(), describe(), record() (+8 more)

### Community 64 - "policyEditor"
Cohesion: 0.13
Nodes (13): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), netWarnings(), newNetGuardTab() (+5 more)

### Community 65 - "unix_connect.c"
Cohesion: 0.40
Nodes (13): attack_connect(), attack_dgram(), attack_uring(), denied(), fill_addr(), handoff(), hold(), main() (+5 more)

### Community 66 - "wm.go"
Cohesion: 0.11
Nodes (25): detectWM(), reachable(), centered(), Detect(), FitSize(), Session, hyprlandSocket(), i3RuntimeSocket() (+17 more)

### Community 68 - "linkView"
Cohesion: 0.43
Nodes (3): linkTarget(), entryKind, linkView

### Community 69 - "NetGuardEvent"
Cohesion: 0.13
Nodes (14): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), snapNetGuard, Mode (+6 more)

### Community 70 - "Vault"
Cohesion: 0.07
Nodes (41): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), encryptDirectories(), verifyEncryptionState(), collectFilesystemPrereqs(), confirmRunPrereq(), resolveFilesystemPrereqs() (+33 more)

### Community 71 - "daemon_multicall_test.go"
Cohesion: 0.17
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 72 - "journal"
Cohesion: 0.14
Nodes (6): fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal

### Community 73 - "fcntl"
Cohesion: 0.10
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 74 - "BinaryEntry"
Cohesion: 0.08
Nodes (22): binaryVerifyState, HashStore, keptEntries, sharedHashEntry, BinaryEntry, checkBtrfsKeys(), ensureBtrfsLayout(), canonicalBinaryPath() (+14 more)

### Community 75 - "buildTrustedSet"
Cohesion: 0.08
Nodes (33): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), allowedEnv(), binaryLaunchRules(), launchRulesByPath(), TestAllowedEnvKeepsOnlyRootPlaced(), TestLaunchRulesByPath(), buildTrustedSet() (+25 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "runMonitor"
Cohesion: 0.08
Nodes (26): DrainUntilSignal(), BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag() (+18 more)

### Community 78 - "listFacultatives"
Cohesion: 0.17
Nodes (17): installFacultatives(), installSSHAgentFacultative(), listFacultatives(), runFacultatives(), facultativeUser(), TestInstallSSHAgentFacultative(), TestListFacultativesIgnoresDaemonConf(), TestListFacultativesMarksGuardedSSH() (+9 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.21
Nodes (11): copy(), is_verb(), main(), read_once(), usage(), dump_memory(), find_all_rw_regions(), find_heap() (+3 more)

### Community 80 - "GuardInodeKey"
Cohesion: 0.03
Nodes (83): scriptIdentity(), confirmUserPlaced(), forwardCheck(), refusedReplacements, VettedInode, Multicall, MulticallKind, confinedEntry() (+75 more)

### Community 81 - "recordFlow"
Cohesion: 0.10
Nodes (25): checkAndApply(), putConfig(), encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), offerBackupDeletion() (+17 more)

### Community 82 - "sync.Mutex"
Cohesion: 0.11
Nodes (11): newActivity(), newExeResolver(), processKey(), guardRow(), TestRingKeepsNewest(), Counts, newRing(), activity (+3 more)

### Community 83 - "rcUser"
Cohesion: 0.11
Nodes (44): setupBunTmpdir(), rcTarget, BunTmpDir(), TestHasRCBlocksTrackEnsureAndRemove(), TestRemoveSSHAgentShellEnvKeepsSSHConfig(), applyRC(), bunTargets(), EnsureAddKeysToAgent() (+36 more)

### Community 84 - "OpenSystemPlaced"
Cohesion: 0.11
Nodes (25): placedWalk, treeJudge, confineBelow(), inUserHome(), vouchedFd(), GlobSystemPlaced(), OpenSystemPlaced(), placedDir() (+17 more)

### Community 85 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 86 - "exe_supersede.h"
Cohesion: 0.21
Nodes (15): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+7 more)

### Community 87 - "loginServer"
Cohesion: 0.07
Nodes (24): markReady(), markReadyWithLogin(), notifySystemdReady(), hasSentinel(), makeRuntimeDir(), openUserRuntimeDir(), plantSentinel(), runtimePathsFor() (+16 more)

### Community 88 - "vaultfs.go"
Cohesion: 0.09
Nodes (26): entryOf(), DefaultNewFileMode(), FileModeToUnixMode(), Vault, IsBinaryContent(), IsOctalMode(), Open(), ParseMode() (+18 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "DumpHook"
Cohesion: 0.32
Nodes (4): formatEntry(), sortedKeys(), sortStrings(), DumpHook

### Community 92 - "tempGrantSetup"
Cohesion: 0.23
Nodes (9): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), fakeTempGrant (+1 more)

### Community 93 - "netinject.c"
Cohesion: 0.12
Nodes (15): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+7 more)

### Community 95 - "NetGuard"
Cohesion: 0.08
Nodes (14): lsmHook, trustHook, foreignPinnedLink(), PumpRingbuf(), reasonText(), statInodeKey(), judgePasses(), parseProcs() (+6 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "time.Duration"
Cohesion: 0.14
Nodes (14): procStats, statsSample, formatAgo(), cpuTime(), humanBytes(), rssBytes(), New(), recv() (+6 more)

### Community 100 - "gui/session.go"
Cohesion: 0.20
Nodes (15): ancestorEnv(), invokingHome(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField(), TestDisplayValue() (+7 more)

### Community 101 - "hashCache"
Cohesion: 0.13
Nodes (17): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, BinaryStat, DurableStat, HashCache (+9 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "writeWithin"
Cohesion: 0.25
Nodes (8): atomicWriteAt(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync(), writeWithin()

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "ptrace_seize.c"
Cohesion: 0.20
Nodes (11): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file(), main(), peek_buffer(), run_tracer() (+3 more)

### Community 106 - "recordops.c"
Cohesion: 0.36
Nodes (15): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+7 more)

### Community 107 - "Open"
Cohesion: 0.10
Nodes (31): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+23 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 110 - "bufio.Reader"
Cohesion: 0.17
Nodes (14): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readBareRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest() (+6 more)

### Community 111 - "ConfinedEntry"
Cohesion: 0.29
Nodes (19): ConfinedEntry(), HashStats(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash(), TestPersistentStoreFilledAfterSyncAndHit() (+11 more)

### Community 112 - "UI"
Cohesion: 0.07
Nodes (29): confirmKeyDeletion(), run(), runUninstall(), RunWith(), removeKeyAndEmptyDir(), removeMasterKey(), TestRemoveKeyAndEmptyDir(), TestRunWithRefusesConcurrentRun() (+21 more)

### Community 113 - "KernelDev"
Cohesion: 0.11
Nodes (23): fdStat, physicalParent(), statFD(), KernelDev(), BtrfsDevices(), btrfsSysfsDevices(), FirstOnBtrfs(), mountinfoDev() (+15 more)

### Community 115 - "maintainerKeeper"
Cohesion: 0.18
Nodes (12): appRunning(), libDirSig(), maintainedApps(), ownedBy(), sameDevice(), startMaintainers(), TestMaintainedAppsNeedTheUsersReadOnlyLibDir(), maintainedApp (+4 more)

### Community 116 - "PlanDaemonUnitUpdate"
Cohesion: 0.15
Nodes (14): TestDaemonUnitWithMetadataOutput(), DaemonUnitWithMetadataOutput(), ConfirmUnitFunc, InstalledDaemonUnit(), PlanDaemonUnitUpdate(), ShippedDaemonUnit(), fakeBinary(), TestPlanDaemonUnitUpdate() (+6 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - ".netIntegrityContainer"
Cohesion: 0.28
Nodes (5): IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 121 - "Gate"
Cohesion: 0.14
Nodes (6): placedKey(), Caller, Gate, Role, Listen(), Judge

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "inet_connect.c"
Cohesion: 0.26
Nodes (15): attack_connect(), attack_fastopen(), attack_uring(), denied(), dial(), fill_addr(), handoff(), listen_on() (+7 more)

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 127 - "classify"
Cohesion: 0.22
Nodes (12): call(), TestClassify(), TestClassifyCreateItemReplace(), TestClassifyOpenCallsExact(), arg(), attributesOf(), classify(), classifyCollection() (+4 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "newTestTarget"
Cohesion: 0.15
Nodes (19): Target, MonitorRepository, FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard() (+11 more)

### Community 130 - "catalogRefresher"
Cohesion: 0.10
Nodes (18): drain(), newCatalogRefresher(), readEvents(), refreshAdmits(), swapFile(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), addSystemBinary(), catalogPatterns() (+10 more)

### Community 131 - "planUpdaters"
Cohesion: 0.14
Nodes (18): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+10 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "Finder"
Cohesion: 0.13
Nodes (19): findButton, span, findAll(), Finder, lowerRunes(), NewFinder(), newFindInput(), overlay() (+11 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "Daemon configuration"
Cohesion: 0.14
Nodes (13): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Launch rules: `[launch_rules]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Port sections, Process inspectors (+5 more)

### Community 143 - "serve.go"
Cohesion: 0.07
Nodes (32): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+24 more)

### Community 144 - "liveSession"
Cohesion: 0.15
Nodes (16): dialLiveSession(), parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), eventsLabel(), forwardBinaries(), holdForward() (+8 more)

### Community 145 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.17
Nodes (22): class, pos, classOf(), column(), dedent(), Edit(), leadingIndent(), lines() (+14 more)

### Community 146 - "Ledger"
Cohesion: 0.11
Nodes (5): applySchema(), Ledger, Pending, Source, fakeLedger

### Community 147 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.13
Nodes (15): App maintainers (Steam runtime updates), Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials (+7 more)

### Community 148 - "fakeGuardRepo"
Cohesion: 0.07
Nodes (16): binaryStatOf(), TempBinary, TemporaryGrant, TempRule, applyTempGrants(), closeTempBinaries(), TemporaryAccess, logTempGrant() (+8 more)

### Community 149 - "os.FileInfo"
Cohesion: 0.21
Nodes (5): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata()

### Community 150 - "find"
Cohesion: 0.22
Nodes (16): deletePostBackups(), listBackups(), restoreBackups(), offerBackupCleanup(), TestMovedFlagError(), movedFlagError(), Delete(), Find() (+8 more)

### Community 152 - "EventType"
Cohesion: 0.08
Nodes (17): selfEventsFor(), snapEngines, snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, planTempBlock() (+9 more)

### Community 154 - "buildOneGuard"
Cohesion: 0.12
Nodes (25): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), startSelfGuard(), startSelfGuards(), selfGuards (+17 more)

### Community 155 - "watchSet"
Cohesion: 0.19
Nodes (8): isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet, TestNearestRootClaims(), nearestRootClaims()

### Community 156 - "copyRegularAt"
Cohesion: 0.31
Nodes (10): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), fchownFromInfo(), futimesFromInfo(), readDirNames() (+2 more)

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 159 - "backend_test.go"
Cohesion: 0.23
Nodes (22): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+14 more)

### Community 161 - ".Layout"
Cohesion: 0.44
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 162 - "NetEventType"
Cohesion: 0.05
Nodes (26): snapNetMonitor, NetEvent, NetEventType, NetRecord, eventsOrAll(), BinariesSummary(), TestBinariesSummary(), FormatAddr() (+18 more)

### Community 164 - "promptPassword"
Cohesion: 0.22
Nodes (11): applyEdit(), editAgain(), fetchConfig(), nextConfig(), runEditConfig(), promptPassword(), promptPasswordWith(), PromptNewPassword() (+3 more)

### Community 165 - "feed.go"
Cohesion: 0.07
Nodes (42): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+34 more)

### Community 166 - "shell"
Cohesion: 0.05
Nodes (43): daemonPage, homePage, pageDef, pageID, reel, resizeWatch, runHost, shell (+35 more)

### Community 167 - "sources_test.go"
Cohesion: 0.16
Nodes (11): rootOwnedInode(), rootOwnedSafe(), SystemTrusted(), systemTool(), archiveMembers(), decompressed(), newSources(), tarMembers() (+3 more)

### Community 168 - "RunWith"
Cohesion: 0.13
Nodes (19): installGUI(), needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), stagedBinary, Mode (+11 more)

### Community 169 - "Run"
Cohesion: 0.33
Nodes (8): planFrom(), newStage(), Run(), runSandbox(), tail(), TestCopyScriptHoldsOnlyTheConfirmedContent(), Spec, stage

### Community 170 - "SelectOrphans"
Cohesion: 0.20
Nodes (17): modifiedContextWithSource(), appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors() (+9 more)

### Community 172 - "Limitations"
Cohesion: 0.18
Nodes (11): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Secrets served over IPC (+3 more)

### Community 173 - "runOfflineEdit"
Cohesion: 0.09
Nodes (22): auditAfterEdit(), auditEntry(), auditFindings(), auditFindingsWithConfig(), auditTree(), coveredByWatch(), relOrPath(), TestAuditTree() (+14 more)

### Community 175 - "TrustGuard"
Cohesion: 0.05
Nodes (24): argvMatches, statFunc, tempGrant, trustRows, setSupersedeTrust(), liveTrust(), loadPinnedExeMaps(), stripPinnedAllow() (+16 more)

### Community 176 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 178 - "update/update.go"
Cohesion: 0.11
Nodes (42): installHooks(), staged(), applyUpdate(), archAssets(), assetsFor(), channelOf(), checkUpdatePreconditions(), confirmUnit() (+34 more)

### Community 179 - "selection_test.go"
Cohesion: 0.25
Nodes (10): clipboardText(), setClipboard(), cellOf(), repeat(), selKey(), TestSelectionEditing(), TestSelectionLargeDelete(), TestSelectionMouseDrag() (+2 more)

### Community 180 - "IntegrationSuite"
Cohesion: 0.24
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 181 - "How it works"
Cohesion: 0.22
Nodes (9): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener (+1 more)

### Community 182 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 186 - "Binary trust"
Cohesion: 0.22
Nodes (9): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch rules, Launch scan, Multicall binaries, The dynamic loader (+1 more)

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "github.com/spf13/cobra.Command"
Cohesion: 0.09
Nodes (31): credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials(), TestSetupDumpLogRefusesSymlink() (+23 more)

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.40
Nodes (5): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog, What an entry can't cover

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network monitor: see what a program talks to

### Community 195 - "openBunTmp"
Cohesion: 0.27
Nodes (11): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir() (+3 more)

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 199 - "filevault_test.go"
Cohesion: 0.09
Nodes (45): harnessSuite, classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord(), newFileVaultAEAD() (+37 more)

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "pinstate.go"
Cohesion: 0.18
Nodes (22): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+14 more)

### Community 202 - "DaemonEvent"
Cohesion: 0.14
Nodes (24): ServeConfig, catchLifecycleSignals(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI(), runTUI() (+16 more)

### Community 203 - "WithBottomBar"
Cohesion: 0.31
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 204 - "runBackend"
Cohesion: 0.24
Nodes (9): reserveStdout(), runBackend(), LiveModeAvailable(), requireRoot(), runEditProtected(), runLive(), anyDaemonRunning(), DaemonRunning() (+1 more)

### Community 205 - "batchInodeRows"
Cohesion: 0.27
Nodes (9): inodeBatchFn, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows(), walkInodeBatches(), countInodeRows() (+1 more)

### Community 206 - "monitor"
Cohesion: 0.07
Nodes (30): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks (+22 more)

### Community 209 - "revertSystemFiles"
Cohesion: 0.24
Nodes (10): removeSessionHooks(), revertSystemFiles(), RemoveLoginDropIn(), removeFile(), SyncSecretServiceFiles(), LegacyCatalogRefreshFiles(), RemoveBinSymlink(), RemoveBinSymlinkAt() (+2 more)

### Community 210 - "main_test.go"
Cohesion: 0.20
Nodes (5): infraContainerPaths(), TestIntegrationSuite(), TestMain(), helperChild(), TestMain()

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 212 - "BusClient"
Cohesion: 0.31
Nodes (7): BusClient, AttestBusClient(), busClientKey(), execStamp(), imageOf(), pidfdPid(), procIdentity()

### Community 213 - "procstats.go"
Cohesion: 0.38
Nodes (8): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic()

### Community 214 - "writeEvent"
Cohesion: 0.33
Nodes (8): UIDResolver, NewUIDResolver(), writeEvent(), allowedEvent(), deniedEvent(), TestWriteEventAllowed(), TestWriteEventAllowedSuppressed(), TestWriteEventBlocked()

### Community 216 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 217 - "SecretStoreNotes"
Cohesion: 0.39
Nodes (7): dockerKeyringHelper(), gitCredentialHelpers(), readOwnedRegular(), SecretStoreNotes(), TestReadOwnedRegular_RefusesSymlinkAndOtherOwner(), TestSecretStoreNotes(), WarnSecretStores()

### Community 218 - "parseDaemonEvents"
Cohesion: 0.15
Nodes (6): IntegrationSuite, IntegrationSuite, IntegrationSuite, parseDaemonEvents(), daemonEvent, IntegrationSuite

### Community 220 - "WithSuspectsFrom"
Cohesion: 0.36
Nodes (4): WithSuspectsFrom(), RealNetEngines(), netguardHarness, Option

### Community 221 - "maintainer/plan.go"
Cohesion: 0.31
Nodes (7): isWhiteout(), opaque(), pathByte(), pathRecord(), TestPathRecord(), opKind, walker

### Community 222 - "ValidateTargets"
Cohesion: 0.19
Nodes (9): existingSection(), guardedMatch(), overlaps(), refusedPath(), sectionFor(), validateTarget(), ValidateTargets(), Existing (+1 more)

### Community 223 - "ReadInstallCert"
Cohesion: 0.36
Nodes (8): InstallCert, ReadInstallCert(), RemoveInstallCert(), RequireNoInstallCert(), TestInstallCertRoundTrip(), TestInstallCertSymlinkRefuses(), TestInstallCertUnreadableRefuses(), withCertPath()

### Community 224 - "downloadFile"
Cohesion: 0.24
Nodes (7): downloadFile(), downloadReleaseFiles(), fetchRelease(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleaseByTag(), progressReader

### Community 225 - "supersede_probe.c"
Cohesion: 0.76
Nodes (6): answer(), main(), read_ok(), read_secret(), run_child(), spin()

### Community 226 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.04
Nodes (54): recordTickCmd(), runNetworkMonitor(), runTUI(), recordModel, diffModel, editorModel, listenForDaemonEvents(), NewDaemonModel() (+46 more)

### Community 227 - "AddServeFlags"
Cohesion: 0.43
Nodes (7): AddHeadlessFlag(), AddServeFlags(), init(), init(), init(), init(), init()

### Community 228 - "mountTable"
Cohesion: 0.13
Nodes (10): MountInfo, mountinfoUnitTest, MountByID(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo(), userMountable() (+2 more)

### Community 231 - "fileedit_symlink_test.go"
Cohesion: 0.36
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 232 - "TestReadFrame"
Cohesion: 0.40
Nodes (4): TestReadFrame(), decodeFrame(), encodeWithSerial(), header()

### Community 233 - "DiscoverInfraBinaries"
Cohesion: 0.33
Nodes (6): DiscoverInfraBinaries(), infraFromRunning(), runningExecutables(), pathIsRunning(), TestDiscoverInfraBinaries(), TestInfraFromRunning()

### Community 235 - "inode_dev.h"
Cohesion: 0.60
Nodes (3): ctx_ptr(), sb_dev(), guard_inode_free()

### Community 236 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 237 - "VouchedDevs"
Cohesion: 0.29
Nodes (7): hostVouchedDevs(), isFuseType(), parseMountDev(), VouchedDevs(), TestParseMountDevMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed()

### Community 238 - "GenerateConf"
Cohesion: 0.05
Nodes (67): logRefreshChange(), TestLogRefreshChange(), discordCatalogEntry(), groupedDiscordConf(), TestGroupedSectionAddressedByEncryptionRoot(), TestPatchCatalogSectionGroupedConfig(), TestSectionsFromCandidatesSkipsMissingWatchSubPaths(), TestSetSectionWhitelistPreservesGroupStructure() (+59 more)

### Community 239 - "parseReservedPorts"
Cohesion: 0.50
Nodes (3): parseReservedPorts(), readReservedPorts(), TestParseReservedPorts()

### Community 240 - "shQuote"
Cohesion: 0.17
Nodes (5): IntegrationSuite, IntegrationSuite, overlayPlanter(), shQuote(), IntegrationSuite

### Community 242 - "User"
Cohesion: 0.05
Nodes (59): askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchers(), bunLaunchersFor(), gatherPerUserSetup(), addedResources(), applyDiffAdditions(), TestAddedResourcesSkipsAlreadyGuardedSSH() (+51 more)

### Community 243 - "parseKernelRelease"
Cohesion: 0.67
Nodes (4): ciliumVersion, parseKernelRelease(), TestParseKernelRelease(), unameKernelVersion()

### Community 244 - "Check"
Cohesion: 0.50
Nodes (4): Check(), collectDiagnostics(), readSysctl(), requireKernel51()

### Community 249 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 256 - "plan_test.go"
Cohesion: 0.40
Nodes (11): scanLib(), fixture(), TestApplyRefusesAStageChangedAfterValidation(), TestApplyWritesValidatedContent(), TestScanAndValidateAcceptsLibDirContent(), TestScanRefusesSpecialFiles(), TestScanUntouchedLibDir(), TestTouches() (+3 more)

## Knowledge Gaps
- **33 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener/tools/codecheck`, `open_by_handle_at POC` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 664 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **45 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `daemon_multicall_test.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `main_test.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `.netIntegrityContainer`, `IntegrationSuite`, `parseDaemonEvents`, `.TestRecord_GuardParity`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `go_pkg_strings` be split into smaller, more focused modules?**
  _Cohesion score 0.11002518695724052 - nodes in this community are weakly interconnected._
- **Why does `Guard` connect `Guard` to `go_pkg_strings`, `time.Duration`, `engine`, `guardUnitTest`, `DaemonEvent`, `BinaryEntry`, `GuardEvent`, `batchInodeRows`, `binaryVetter`, `TrustGuard`, `GuardInodeKey`, `sync.Mutex`, `Resource`, `loginServer`, `buildOneGuard`, `time.Time`?**
  _High betweenness centrality (0.027) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.018337651983256927 - nodes in this community are weakly interconnected._
- **Why does `GuardRepository` connect `Resource` to `go_pkg_strings`, `GuardEvent`, `fakeGuardRepo`, `usecase/daemon_test.go`, `TestMonitorUseCaseLifecycle`?**
  _High betweenness centrality (0.014) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.09441441441441441 - nodes in this community are weakly interconnected._