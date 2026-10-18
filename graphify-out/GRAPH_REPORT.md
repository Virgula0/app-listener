# Graph Report - app-listener  (2026-10-18)

## Corpus Check
- 471 files · ~1,690,287 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 6162 nodes · 26191 edges · 227 communities (193 shown, 34 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2459 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `2b586627`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- Vault
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- rcUser
- guardUnitTest
- mkdirs
- find_test.go
- policyEditor
- protectedPage
- github.com/charmbracelet/bubbletea.Model
- .launch
- os.FileMode
- IntegrationSuite
- guard_trust.bpf.c
- IntegrationSuite
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth_test.go
- TrustGuard
- guard.bpf.c
- NetEventType
- time.Time
- .FilterExistingWhitelist
- controlServer
- UI
- serve.go
- IntegrationSuite
- procstats.go
- picker
- github.com/spf13/cobra.Command
- IntegrationSuite
- backend
- editControlSession
- fileedit_model_test.go
- os.File
- IntegrationSuite
- absPath
- monitorUnitTest
- BinaryEntry
- GuardEvent
- render.go
- stdio
- editView
- fadeScroll
- LibraryClosure
- GenerateConf
- proto.go
- daemonconfig.go
- tabs.go
- Open
- run
- icon
- GuardInodeKey
- filevault_test.go
- update_test.go
- newTestTarget
- exe_supersede.h
- Network
- decryptStep
- bpfstats/main.go
- gui/session.go
- runMonitor
- filevault.go
- highlight_test.go
- fcntl
- dialLiveSession
- Guard
- FileTarget
- btrfsLayoutFrom
- lockAndDeprovision
- Hash
- RefreshSection
- CopyTreeWithProgress
- process_vm_readv.c
- Monitor
- journal
- planUpdaters
- WithBottomBar
- compileGlobNames
- Config
- github.com/charmbracelet/bubbles/textarea.Model
- netModel
- liveSession
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- DaemonEvent
- changelog.go
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- catalogRefresher
- fileedit_symlink_test.go
- deploy
- Well-known bypass classes
- KernelDev
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- tempGrantSetup
- Build & sign release assets (reusable workflow)
- installPage
- Active
- mountsUnitTest
- NetworkMonitorRepository
- buildTrustedSet
- IntegrationSuite
- IntegrationSuite
- net_tester/main.go
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- CandidateDir
- IntegrationSuite
- WebSocket /ws client
- launch_probe.c
- install.sh
- mc_tag.h
- PlanDaemonUnitUpdate
- make test-integration (rootful Docker suite)
- RunWith
- supersede_probe.c
- cursor.go
- trace-app-libs.sh
- IntegrationSuite
- run
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- TempBinary
- time.Duration
- IntegrationSuite
- runOfflineEdit
- changelog_test.go
- RawTarget
- Execute
- SelectOrphans
- fetchReleases
- NewEventFanout
- ConfinedEntry
- EventType
- layout.go
- copyRegularAt
- auditGroupCoverage
- Installation
- backend_test.go
- graphify-refresh.sh
- vaultfs.go
- inodeChain
- configedit_test.go
- pinstate.go
- feed.go
- shell
- runView
- controlManager
- runGenKey
- Finder
- eventUnitTest
- Limitations
- Troubleshooting
- IntegrationSuite
- resolveFilesystemPrereqs
- Daemon configuration
- .judgeProcess
- downloadAndVerify
- RunCmd
- Compatibility
- How it works
- fakeDaemon
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- NetworkMonitor
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- fakeVault
- harnessSuite
- edit-protected
- monitor: see what touches a path
- kernel-prof.sh
- Files
- pathCache
- ReadInstallCert
- writeWithin
- bpfLsmListed
- NetGuardEvent
- parseGuardEvent
- Desktop GUI
- IntegrationSuite
- Ledger
- inode_dev.h
- .eventTypeAllowed
- WithSuspectsFrom
- DiscoverInfraBinaries
- ValidatePassword
- CleanupStalePins
- sync.Mutex
- Vault
- checkRunningBinaryMatchesInstalled
- golang.org/x/sys/unix.Timespec
- bleedLayout
- TestMain
- MulticallError
- launchKey
- partitionEncryptionRoots

## God Nodes (most connected - your core abstractions)
1. `Guard` - 117 edges
2. `GuardInodeKey` - 96 edges
3. `Config` - 95 edges
4. `Load()` - 89 edges
5. `shell` - 87 edges
6. `IntegrationSuite` - 82 edges
7. `absPath()` - 81 edges
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

## Communities (227 total, 34 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.09
Nodes (16): tempJournalRow, invalidPassword, launchRule, trustEvent, textAnswer, btrfsIoctlFsInfoArgs, infraContainerPaths(), TestIntegrationSuite() (+8 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (136): plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestDiffCatalogFlagRegistered(), openTemp(), TestAuthSucceededResetsBudget(), TestBeginAuthAttemptLocksAfterBudget() (+128 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (9): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig() (+1 more)

### Community 3 - "Vault"
Cohesion: 0.15
Nodes (17): encryptDirectories(), secureResources(), verifyEncryptionState(), askEncryption(), isFileVaultCiphertext(), isRegularFileTarget(), Vault, hasEncryptionPolicy() (+9 more)

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
Cohesion: 0.13
Nodes (34): setupBunTmpdir(), rcTarget, BunTmpDir(), bunTargets(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv(), EnsureBunTmpDir(), ensureRCBlock() (+26 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (19): addErrKind, ReplacementCheck, classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest, runTool() (+11 more)

### Community 9 - "mkdirs"
Cohesion: 0.11
Nodes (42): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), newWatchPattern(), systemBinaryPatterns(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped() (+34 more)

### Community 10 - "find_test.go"
Cohesion: 0.17
Nodes (20): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+12 more)

### Community 11 - "policyEditor"
Cohesion: 0.13
Nodes (10): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), netWarnings(), confirmWarnings() (+2 more)

### Community 12 - "protectedPage"
Cohesion: 0.06
Nodes (20): backendLink, configView, forwardForm, forwardView, installUI, protectedPage, Event, scrollable() (+12 more)

### Community 13 - "github.com/charmbracelet/bubbletea.Model"
Cohesion: 0.12
Nodes (7): renderBullets(), TestNoticeViewFitsAndPlacesOKLast(), TestRenderBulletsVerbatimAndIndented(), noticeModel, servedModel, serveTestModel, sessionModel

### Community 14 - ".launch"
Cohesion: 0.27
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 15 - "os.FileMode"
Cohesion: 0.09
Nodes (30): dupFD(), installConfig(), installFile(), installFileAs(), installSSHAgent(), installSSHAgentEnv(), installUnit(), readAtNoFollow() (+22 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (9): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent (+1 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (58): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+50 more)

### Community 18 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (12): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, dirEntry (+4 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.20
Nodes (44): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+36 more)

### Community 22 - "Resource"
Cohesion: 0.08
Nodes (10): ownerPaths(), Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, uniqueEncryptionRoots(), unlockRoots() (+2 more)

### Community 23 - "auth_test.go"
Cohesion: 0.28
Nodes (19): HashConfigured(), LoadHash(), MigrateLegacyHash(), mustHash(), TestCheckPasswordSharesGuessBudget(), TestMigrateLegacyHashEmptyPlaceholder(), TestMigrateLegacyHashImports(), TestMigrateLegacyHashNeverReplaces() (+11 more)

### Community 24 - "TrustGuard"
Cohesion: 0.05
Nodes (23): exeRow, tempGrant, tempMaskOp, tempRow, trustHook, trustRows, guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent() (+15 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (72): inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), current_is_inspector() (+64 more)

### Community 26 - "NetEventType"
Cohesion: 0.10
Nodes (16): snapNetMonitor, NetBpfEvent, NetEvent, NetEventType, eventsOrAll(), newNetEventChecks(), FormatAddr(), NetEventTypes() (+8 more)

### Community 27 - "time.Time"
Cohesion: 0.13
Nodes (10): LockedError, isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet, Process inspectors (+2 more)

### Community 28 - ".FilterExistingWhitelist"
Cohesion: 0.10
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 29 - "controlServer"
Cohesion: 0.23
Nodes (4): startControlServer(), TestRecordFailureLockout(), configEditor, controlServer

### Community 30 - "UI"
Cohesion: 0.10
Nodes (18): InstallMode, InstallRequest, Options, UninstallHooks, UninstallRequest, UpdateHooks, InstallHooks, StagedRelease (+10 more)

### Community 31 - "serve.go"
Cohesion: 0.08
Nodes (30): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+22 more)

### Community 32 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): guardBinaryFlag(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 33 - "procstats.go"
Cohesion: 0.07
Nodes (28): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+20 more)

### Community 34 - "picker"
Cohesion: 0.08
Nodes (30): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, openPicker(), pickerTree() (+22 more)

### Community 35 - "github.com/spf13/cobra.Command"
Cohesion: 0.07
Nodes (34): AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials() (+26 more)

### Community 37 - "backend"
Cohesion: 0.12
Nodes (8): checkSecret(), backend, backendOp, backendState, fakeLink, Message, Request, RequireDaemonStopped()

### Community 38 - "editControlSession"
Cohesion: 0.13
Nodes (18): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readBareRequest(), readClientLines(), readForwardRequest(), readGrantRequest(), readPathLines() (+10 more)

### Community 39 - "fileedit_model_test.go"
Cohesion: 0.20
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 40 - "os.File"
Cohesion: 0.05
Nodes (60): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), setConfinedHomes() (+52 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.07
Nodes (15): IntegrationSuite, decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest (+7 more)

### Community 42 - "absPath"
Cohesion: 0.08
Nodes (5): shQuote(), IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 43 - "monitorUnitTest"
Cohesion: 0.14
Nodes (3): dirTarget(), newMountTable(), monitorUnitTest

### Community 44 - "BinaryEntry"
Cohesion: 0.10
Nodes (20): binaryVerifyState, HashStore, sharedHashEntry, BinaryEntry, BinariesSummary(), claimHash(), hashFileShared(), hashPathShared() (+12 more)

### Community 45 - "GuardEvent"
Cohesion: 0.04
Nodes (26): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), holdHeadless(), HeadlessLine(), resolveGuardConfig(), runGuard(), runGuardHeadless() (+18 more)

### Community 46 - "render.go"
Cohesion: 0.11
Nodes (13): lineStyles, viewportState, baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces(), setGutter() (+5 more)

### Community 47 - "stdio"
Cohesion: 0.07
Nodes (4): put(), wait_then_read(), hold_forever(), main()

### Community 48 - "editView"
Cohesion: 0.11
Nodes (9): codeEntry, editView, handButton, Entry, selectable(), fileSize(), joinRel(), newCodeEntry() (+1 more)

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (26): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+18 more)

### Community 50 - "LibraryClosure"
Cohesion: 0.07
Nodes (29): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), isScript(), ldConfInclude(), ldSoConfDirs() (+21 more)

### Community 51 - "GenerateConf"
Cohesion: 0.06
Nodes (53): TestDiffMergePreservesExistingAndParses(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestGeneratedLibraryBlockRoundTrips(), TestGroupedSectionAddressedByEncryptionRoot(), TestSectionsFromCandidatesSkipsMissingWatchSubPaths() (+45 more)

### Community 52 - "proto.go"
Cohesion: 0.06
Nodes (22): Client, File, User, launch(), TestCleanRel(), TestHelperBackend(), TestReaderCapsLines(), CleanRel() (+14 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.14
Nodes (37): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+29 more)

### Community 54 - "tabs.go"
Cohesion: 0.03
Nodes (23): appTheme, binaryRow, clickPanel, iconTile, iconTileRenderer, navItem, navItemRenderer, panelRenderer (+15 more)

### Community 55 - "Open"
Cohesion: 0.18
Nodes (19): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+11 more)

### Community 56 - "run"
Cohesion: 0.10
Nodes (30): parseMaintenanceFlags(), runMaintenanceMode(), TestValidateMaintenanceFlags(), validateMaintenanceFlags(), deletePostBackups(), listBackups(), restoreBackups(), confirmKeyDeletion() (+22 more)

### Community 57 - "icon"
Cohesion: 0.08
Nodes (37): checkGrid, filesPage, homePage, multiAnswer, panel, statusTile, newBinaryList(), newBinaryRow() (+29 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.04
Nodes (47): engine, resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), TestSetInspectorsBounded(), appendInodeRow(), batchInodeRows() (+39 more)

### Community 59 - "filevault_test.go"
Cohesion: 0.10
Nodes (39): decryptDirectories(), TestDecryptDirectoriesEmpty(), looksLikeFileVaultRecord(), newFileVaultAEAD(), openFileVault(), sealFileVault(), inode(), mustSubkey() (+31 more)

### Community 60 - "update_test.go"
Cohesion: 0.12
Nodes (15): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFilterChannel(), testKeyPair(), TestNewerThanInstalled(), TestNewerThanStable() (+7 more)

### Community 61 - "newTestTarget"
Cohesion: 0.23
Nodes (16): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+8 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.21
Nodes (15): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+7 more)

### Community 63 - "Network"
Cohesion: 0.10
Nodes (10): binaryList, netGuardPill(), newExeResolver(), guardRow(), clonePolicy(), NetExtra, NetPolicy, Network (+2 more)

### Community 64 - "decryptStep"
Cohesion: 0.29
Nodes (7): askDecrypt(), cleanOrphanedMetadata(), decryptStep(), pickDirsToDecrypt(), CleanOrphanedMetadata(), LivePolicyDescriptors(), TestLivePolicyDescriptorsTmpfs()

### Community 65 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 66 - "gui/session.go"
Cohesion: 0.07
Nodes (41): ancestorEnv(), cleanupSessionEnv(), invokingHome(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField() (+33 more)

### Community 67 - "runMonitor"
Cohesion: 0.10
Nodes (17): CheckEBPF(), MakeDisplayPaths(), ParseNetEventsFlag(), runHeadless(), runMonitor(), runTUI(), runNetworkMonitor(), runTUI() (+9 more)

### Community 68 - "filevault.go"
Cohesion: 0.25
Nodes (12): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), openFileVaultWithMasterKey(), recoverFileInPlace(), rejectSidecarSymlink() (+4 more)

### Community 69 - "highlight_test.go"
Cohesion: 0.13
Nodes (18): fg(), Highlighter, lexerFor(), NewHighlighter(), shebangLexer(), editSession(), sampleGo(), TestDaemonConfTokens() (+10 more)

### Community 70 - "fcntl"
Cohesion: 0.09
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 71 - "dialLiveSession"
Cohesion: 0.13
Nodes (18): auditAfterEdit(), auditFindings(), dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), putConfig() (+10 more)

### Community 72 - "Guard"
Cohesion: 0.05
Nodes (29): backingDeviceUnion(), lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), deferredBinary, rootHandle, openRoot(), TestWithChmodDropWriteNeedsReadOnly() (+21 more)

### Community 73 - "FileTarget"
Cohesion: 0.19
Nodes (5): Target, guardPill(), FileTarget, Policy, FileGuard

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 75 - "lockAndDeprovision"
Cohesion: 0.21
Nodes (8): deprovisionKind, classifyDeprovision(), classifyDeprovisionErr(), applyRawKeyPolicy(), isDeprovisionBusy(), isDeprovisionMissing(), lockAndDeprovision(), modifiedContextWithSource()

### Community 76 - "Hash"
Cohesion: 0.17
Nodes (17): CheckPassword(), Hash(), openLedger(), OriginOf(), parseHash(), StoreHash(), TestHashVerifyRoundTrip(), TestVerifyMalformed() (+9 more)

### Community 77 - "RefreshSection"
Cohesion: 0.11
Nodes (24): refreshCatalogConfig(), TestPatchCatalogSectionGroupedConfig(), vaultScan(), RefreshOptions, SectionScan, TestUnifiedDiff(), TestUnifiedDiffIdentical(), UnifiedDiff() (+16 more)

### Community 78 - "CopyTreeWithProgress"
Cohesion: 0.16
Nodes (8): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata(), CopyTreeWithProgress(), TestCopyTreeProgress(), treeSize()

### Community 79 - "process_vm_readv.c"
Cohesion: 0.18
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "Monitor"
Cohesion: 0.13
Nodes (3): NewMonitor(), Monitor, watchedInode

### Community 81 - "journal"
Cohesion: 0.18
Nodes (5): fakeGuard, fakeNet, fakeNetGuard, fakeNetMonitor, journal

### Community 82 - "planUpdaters"
Cohesion: 0.15
Nodes (17): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+9 more)

### Community 83 - "WithBottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 84 - "compileGlobNames"
Cohesion: 0.17
Nodes (9): GlobChild, statFunc, compileGlobNames(), globText(), resolveBits(), resolveChildren(), statKey(), TestCompileGlobNames_RepeatedPatternOrsBits() (+1 more)

### Community 85 - "Config"
Cohesion: 0.07
Nodes (58): askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), appendSectionsAndEdit(), applyDiffAdditions(), collectDiffAdditions(), configPaths() (+50 more)

### Community 86 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.11
Nodes (34): class, pos, diffModel, editorModel, classOf(), column(), dedent(), Edit() (+26 more)

### Community 87 - "netModel"
Cohesion: 0.21
Nodes (8): formatAddr(), formatNetProto(), formatNetType(), listenForNetEvents(), tickNetStats(), netEventLine, netEventMsg, netModel

### Community 88 - "liveSession"
Cohesion: 0.16
Nodes (13): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), confirmUserPlaced(), displayEntries(), eventsLabel(), forwardBinaries() (+5 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "DaemonEvent"
Cohesion: 0.11
Nodes (33): ServeConfig, UIDResolver, NewUIDResolver(), catchLifecycleSignals(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI() (+25 more)

### Community 92 - "changelog.go"
Cohesion: 0.25
Nodes (17): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), compareStableVersions(), filterChannel(), findUpdate(), latestTime() (+9 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.10
Nodes (10): eventTypeKey(), mcNameKey(), statInodeKey(), judgePasses(), parseProcs(), TestJudgePasses(), TestParseProcs(), hook (+2 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "catalogRefresher"
Cohesion: 0.09
Nodes (20): drain(), logRefreshChange(), newCatalogRefresher(), readEvents(), sameConfig(), swapFile(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads() (+12 more)

### Community 100 - "fileedit_symlink_test.go"
Cohesion: 0.25
Nodes (10): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+2 more)

### Community 101 - "deploy"
Cohesion: 0.16
Nodes (17): applyUpdate(), archAssets(), assetsFor(), checkUpdatePreconditions(), deploy(), deployGUI(), fetchGUI(), guiAssetsFor() (+9 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "KernelDev"
Cohesion: 0.07
Nodes (37): fdStat, BinaryStat, placedWalk, physicalParent(), statFD(), checkBtrfsKeys(), binaryStatOf(), parseMajorMinor() (+29 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "tempGrantSetup"
Cohesion: 0.29
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "installPage"
Cohesion: 0.06
Nodes (20): actionGate, binSource, installAction, installedBinary, installPage, latestRelease, sysStatus, uninstallCard (+12 more)

### Community 110 - "Active"
Cohesion: 0.16
Nodes (12): nextConfig(), showChangelog(), confirmUpdate(), isTerminal(), TestConfirmUpdateRoutesToActiveUI(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), cancelUI, EditText() (+4 more)

### Community 111 - "mountsUnitTest"
Cohesion: 0.12
Nodes (10): MountInfo, mountinfoUnitTest, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo() (+2 more)

### Community 112 - "NetworkMonitorRepository"
Cohesion: 0.13
Nodes (11): NetworkMonitorRepository, NewGuardUseCase(), NetworkMonitorUseCase, NewNetworkMonitorUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle() (+3 more)

### Community 113 - "buildTrustedSet"
Cohesion: 0.19
Nodes (13): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+5 more)

### Community 116 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "CandidateDir"
Cohesion: 0.08
Nodes (24): entryWriters(), reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName() (+16 more)

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
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 127 - "PlanDaemonUnitUpdate"
Cohesion: 0.19
Nodes (13): TestDaemonUnitWithMetadataOutput(), confirmUnit(), DaemonUnitWithMetadataOutput(), ConfirmUnitFunc, InstalledDaemonUnit(), PlanDaemonUnitUpdate(), ShippedDaemonUnit(), fakeBinary() (+5 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "RunWith"
Cohesion: 0.14
Nodes (17): needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), applyLiveRefresh(), runRefreshCatalogOnly(), softenAutomatedRefreshErr() (+9 more)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "cursor.go"
Cohesion: 0.14
Nodes (13): buttonCursor, handCursors(), handCursorsOn(), newButtonCursor(), showConfirm(), showCustomConfirm(), showDialog(), showErrorDialog() (+5 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "run"
Cohesion: 0.18
Nodes (14): channelOf(), fetchRelease(), InstalledGUIVersion(), InstalledVersion(), Latest(), readInstalledVersion(), run(), runUpdate() (+6 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "TempBinary"
Cohesion: 0.09
Nodes (20): TempBinary, TemporaryGrant, TempRule, loadPinnedExeMaps(), planTempBlock(), ResolveTempBinary(), stripPinnedAllow(), StripPinnedTempAllows() (+12 more)

### Community 143 - "time.Duration"
Cohesion: 0.15
Nodes (13): procStats, statsSample, cpuTime(), humanBytes(), rssBytes(), New(), recv(), TestSubDeliversMissedTickLate() (+5 more)

### Community 145 - "runOfflineEdit"
Cohesion: 0.09
Nodes (19): ParseEventsFlag(), reserveStdout(), runBackend(), LiveModeAvailable(), TestValidateFlagsEditConfig(), runEditProtected(), TestPickEncryptedDirSingleSkipsPicker(), validateFlags() (+11 more)

### Community 146 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 147 - "RawTarget"
Cohesion: 0.31
Nodes (9): BuildEBPFTargets(), RawTarget, isSubDir(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets(), WarnIgnoredFlags() (+1 more)

### Community 148 - "Execute"
Cohesion: 0.20
Nodes (9): releaseBuild(), Requested(), Run(), TestRequested(), TestSubcommandOutranksGUI(), Execute(), initLogger(), subcommand() (+1 more)

### Community 149 - "SelectOrphans"
Cohesion: 0.23
Nodes (16): appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors(), SelectOrphans() (+8 more)

### Community 150 - "fetchReleases"
Cohesion: 0.22
Nodes (8): downloadFile(), downloadReleaseFiles(), fetchReleases(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleases(), TestFetchReleasesHTTPError(), progressReader

### Community 152 - "NewEventFanout"
Cohesion: 0.22
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 153 - "ConfinedEntry"
Cohesion: 0.11
Nodes (34): loadDaemonConfig(), markReady(), noteReloadedSettings(), notifySystemdReady(), resolveConfigPath(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), inodeBatchFn (+26 more)

### Community 154 - "EventType"
Cohesion: 0.10
Nodes (13): snapMonitor, EventType, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), specFromForm(), specSummary() (+5 more)

### Community 155 - "layout.go"
Cohesion: 0.43
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 156 - "copyRegularAt"
Cohesion: 0.31
Nodes (10): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), fchownFromInfo(), futimesFromInfo(), readDirNames() (+2 more)

### Community 157 - "auditGroupCoverage"
Cohesion: 0.22
Nodes (11): auditEntry(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths(), relOrPath() (+3 more)

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "backend_test.go"
Cohesion: 0.22
Nodes (19): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+11 more)

### Community 161 - "vaultfs.go"
Cohesion: 0.09
Nodes (26): entryOf(), DefaultNewFileMode(), FileModeToUnixMode(), Vault, IsBinaryContent(), IsOctalMode(), Open(), ParseMode() (+18 more)

### Community 162 - "inodeChain"
Cohesion: 0.26
Nodes (13): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+5 more)

### Community 163 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 164 - "pinstate.go"
Cohesion: 0.12
Nodes (29): configureDaemonLogging(), runLockdown(), ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode() (+21 more)

### Community 165 - "feed.go"
Cohesion: 0.06
Nodes (41): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+33 more)

### Community 166 - "shell"
Cohesion: 0.06
Nodes (39): pageDef, pageID, reel, resizeWatch, shell, sizedDialog, snapEngines, snapNet (+31 more)

### Community 167 - "runView"
Cohesion: 0.07
Nodes (11): logGrid, logHook, logLine, runHost, runView, stageItem, stageState, logRow() (+3 more)

### Community 168 - "controlManager"
Cohesion: 0.29
Nodes (3): newControlManager(), startControlManager(), controlManager

### Community 169 - "runGenKey"
Cohesion: 0.22
Nodes (10): confirmMasterKeyOverwrite(), runGenKey(), ensureInstalledBinary(), ensureMasterKey(), prepareInstallation(), TestEnsureInstalledBinary(), checkKeyLen(), GenerateMasterKey() (+2 more)

### Community 170 - "Finder"
Cohesion: 0.17
Nodes (9): findButton, span, findAll(), Finder, lowerRunes(), newFindInput(), overlay(), runesEqual() (+1 more)

### Community 171 - "eventUnitTest"
Cohesion: 0.27
Nodes (4): BpfEvent, eventUnitTest, Cstr(), DecodeBpfEvent()

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 175 - "resolveFilesystemPrereqs"
Cohesion: 0.18
Nodes (10): collectFilesystemPrereqs(), confirmRunPrereq(), resolveFilesystemPrereqs(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles(), classifySetupError(), TestClassifySetupErrorGeneric(), TestClassifySetupErrorNotSetup() (+2 more)

### Community 176 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Several guarded trees, one vault, Watch sections

### Community 177 - ".judgeProcess"
Cohesion: 0.21
Nodes (10): codeVouched(), isTraced(), loaderEnv(), readMountDevs(), roFor(), TestCodeVouched(), TestLoaderEnv(), TestReadMountDevs() (+2 more)

### Community 178 - "downloadAndVerify"
Cohesion: 0.33
Nodes (10): installHooks(), staged(), downloadAndVerify(), FetchLatest(), fetchPair(), FetchVerified(), Release, sanityCheckBinary() (+2 more)

### Community 179 - "RunCmd"
Cohesion: 0.08
Nodes (36): installGUI(), buildBinaryIfNeeded(), certifyInstall(), installBinaryOnly(), markComplete(), preflightDeployedBinary(), mustCwd(), revertSystemFiles() (+28 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 186 - "Binary trust"
Cohesion: 0.29
Nodes (7): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 189 - "NetworkMonitor"
Cohesion: 0.17
Nodes (5): NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey(), NetworkMonitor

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

### Community 194 - "fakeVault"
Cohesion: 0.42
Nodes (3): TestDaemonUseCaseReloadRefusedDuringShutdown(), TestDaemonUseCaseStopBlocksUntilAllLocked(), fakeVault

### Community 195 - "harnessSuite"
Cohesion: 0.47
Nodes (3): harnessSuite, EnsureSystemSetup(), isKernelAtLeast54()

### Community 196 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "Files"
Cohesion: 0.36
Nodes (5): Files, FileSpec, newFileTarget(), NormalizeSpec(), overlaps()

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "ReadInstallCert"
Cohesion: 0.18
Nodes (14): invokingIDs(), TestWriteUserFileRefusesSymlinks(), writeUserFile(), InstallCert, ReadInstallCert(), RemoveInstallCert(), RequireNoInstallCert(), TestInstallCertRoundTrip() (+6 more)

### Community 202 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 203 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 204 - "NetGuardEvent"
Cohesion: 0.08
Nodes (24): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runHeadless(), runNetworkGuard(), runTUI(), runHeadless() (+16 more)

### Community 205 - "parseGuardEvent"
Cohesion: 0.25
Nodes (9): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels() (+1 more)

### Community 206 - "Desktop GUI"
Cohesion: 0.15
Nodes (13): Actions, Desktop GUI, File and folder pickers, Files, Home, Install, Network, Protected (+5 more)

### Community 209 - "Ledger"
Cohesion: 0.10
Nodes (12): candidates(), confirmable(), confirmAll(), describe(), hashNow(), run(), fyneLogFilter, applySchema() (+4 more)

### Community 210 - "inode_dev.h"
Cohesion: 0.60
Nodes (3): ctx_ptr(), sb_dev(), guard_inode_free()

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 212 - "WithSuspectsFrom"
Cohesion: 0.36
Nodes (4): WithSuspectsFrom(), RealNetEngines(), netguardHarness, Option

### Community 213 - "DiscoverInfraBinaries"
Cohesion: 0.33
Nodes (6): DiscoverInfraBinaries(), infraFromRunning(), runningExecutables(), pathIsRunning(), TestDiscoverInfraBinaries(), TestInfraFromRunning()

### Community 214 - "ValidatePassword"
Cohesion: 0.33
Nodes (7): TestValidatePassword(), ValidatePassword(), promptPassword(), promptPasswordWith(), PromptNewPassword(), promptNewPassword(), PromptNewPasswordUntilValid()

### Community 215 - "CleanupStalePins"
Cohesion: 0.33
Nodes (6): CleanupStalePins(), foreignPinnedLink(), genOfPinFile(), TestCleanupStalePins(), TestCleanupStalePinsMissingBase(), TestGenOfPinFile()

### Community 216 - "sync.Mutex"
Cohesion: 0.03
Nodes (86): CheckBPFLSM(), multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), buildConcurrency() (+78 more)

### Community 218 - "checkRunningBinaryMatchesInstalled"
Cohesion: 0.40
Nodes (5): checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), fileSHA256(), sameContent(), TestCheckRunningBinaryMatchesInstalled()

### Community 222 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

## Knowledge Gaps
- **31 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC`, `statonly stat/statx metadata leak`, `app-listener-gui best-effort build` (+26 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 567 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **34 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Editor keys` connect `edit-protected` to `liveSession`, `run`?**
  _High betweenness centrality (0.046) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC` to the rest of the system?**
  _31 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.09189607072031192 - nodes in this community are weakly interconnected._
- **Why does `edit-protected` connect `edit-protected` to `README.md`?**
  _High betweenness centrality (0.038) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.019887866285835186 - nodes in this community are weakly interconnected._
- **Why does `Delete()` connect `run` to `daemon/daemon.go`, `edit-protected`?**
  _High betweenness centrality (0.034) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.08957708957708958 - nodes in this community are weakly interconnected._