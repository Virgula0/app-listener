# Graph Report - app-listener  (2026-10-18)

## Corpus Check
- 486 files · ~1,711,888 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 6456 nodes · 27854 edges · 220 communities (182 shown, 38 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2589 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `01de7b0f`
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
- IntegrationSuite
- guardUnitTest
- mkdirs
- BinaryEntry
- policyEditor
- protectedPage
- plan.go
- GuardInodeKey
- editView
- IntegrationSuite
- guard_trust.bpf.c
- wipe
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth_test.go
- bpfstats/main.go
- guard.bpf.c
- TrustGuard
- github.com/charmbracelet/bubbles/textarea.Model
- github.com/charmbracelet/bubbletea.Cmd
- shellrc.go
- Config
- TestMonitorUseCaseLifecycle
- UI
- NetEventType
- picker
- fileedit_model_test.go
- IntegrationSuite
- backend
- controlServer
- uninstall/uninstall.go
- ConfinedEntry
- IntegrationSuite
- absPath
- monitorUnitTest
- desktop.go
- IntegrationSuite
- daemon_multicall_test.go
- stdio
- safeio.go
- fadeScroll
- install_ui.go
- cursor.go
- Client
- watchGroup
- tabs.go
- .apply
- prepareDaemonStart
- mutedText
- engine
- os.File
- CandidateDir
- Recorder
- exe_supersede.h
- Network
- Binary
- pinstate.go
- wm.go
- .FilterExistingWhitelist
- recordFlow
- lockAndDeprovision
- fcntl
- TestCleanRel
- Guard
- FileTarget
- KernelDev
- catalogRefresher
- update_test.go
- compileGlobNames
- process_vm_readv.c
- pickfs.go
- newTestTarget
- planUpdaters
- changelog.go
- net_tester/main.go
- .populateWatchInodes
- render.go
- uninstall_system.go
- liveSession
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- lockRootRecovering
- Monitor
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- io.Writer
- highlight_test.go
- downloadFile
- Well-known bypass classes
- installPage
- jit_provenance.c
- inspector_probe.c
- recordops.c
- tempGrantSetup
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- hashCache
- mountTable
- NetGuardEvent
- buildTrustedSet
- IntegrationSuite
- .step
- EventType
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- .TestRecord_GuardParity
- GenerateConf
- WebSocket /ws client
- launch_probe.c
- install.sh
- mc_tag.h
- .filterEvent
- make test-integration (rootful Docker suite)
- .launch
- supersede_probe.c
- install/install.go
- trace-app-libs.sh
- IntegrationSuite
- dialLiveSession
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- GuardEvent
- time.Duration
- auditGroupCoverage
- newTestVetter
- session_test.go
- sync.Mutex
- daemon: permanent protection with encryption at rest
- fscrypt/orphans.go
- watchPattern
- fakeStore
- New
- IntegrationSuite
- buildOneGuard
- layout.go
- IntegrationSuite
- changelog_test.go
- Installation
- backend_test.go
- graphify-refresh.sh
- vaultfs.go
- inodeChain
- runDaemon
- FileEvent
- feed.go
- shell
- runMonitor
- Ledger
- taint_victim.c
- find.go
- DiscoverInfraBinaries
- Limitations
- Troubleshooting
- IntegrationSuite
- IntegrationSuite
- dialog_test.go
- copyRegularAt
- Vault
- fileedit_symlink_test.go
- Compatibility
- How it works
- github.com/spf13/cobra.Command
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- github.com/stretchr/testify/suite.Suite
- writeWithin
- edit-protected
- monitor: see what touches a path
- kernel-prof.sh
- .refreshOnce
- pathCache
- .runEditorHarness
- DaemonEvent
- TestBatchInodeRows
- vouchedDevs
- parsePasswd
- Desktop GUI
- IntegrationSuite
- WithBottomBar
- .eventTypeAllowed
- bleedLayout
- IntegrationSuite
- LibraryClosure
- update/update.go
- Daemon configuration
- time.Time
- bpfLsmListed
- headless.go
- Serve

## God Nodes (most connected - your core abstractions)
1. `Guard` - 118 edges
2. `Config` - 107 edges
3. `GuardInodeKey` - 99 edges
4. `Load()` - 89 edges
5. `shell` - 87 edges
6. `IntegrationSuite` - 82 edges
7. `absPath()` - 82 edges
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
- `resolveAssets()` --calls--> `asset()`  [INFERRED]
  cmd/functions/update/update.go → internal/gui/icons.go

## Import Cycles
- None detected.

## Communities (220 total, 38 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.11
Nodes (21): tempJournalRow, Progress, User, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs (+13 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (134): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+126 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (10): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote() (+2 more)

### Community 3 - "Vault"
Cohesion: 0.20
Nodes (10): isFileVaultCiphertext(), isRegularFileTarget(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr(), newBoundedKeyFn(), readKey(), TestNewBoundedKeyFnFirstCall() (+2 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.06
Nodes (82): TestSetSectionWhitelistPreservesGroupStructure(), Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs() (+74 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.11
Nodes (45): current_exe(), emit_dentry(), emit_event(), emit_event_kern(), emit_file_write(), emit_kern(), emit_user(), fill_identity() (+37 more)

### Community 7 - "IntegrationSuite"
Cohesion: 0.16
Nodes (8): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (18): addErrKind, classifyAddErr(), eventMask(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest, runTool() (+10 more)

### Community 9 - "mkdirs"
Cohesion: 0.10
Nodes (44): fileID(), openBunTmp(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), newWatchPattern(), newFakeInotify() (+36 more)

### Community 10 - "BinaryEntry"
Cohesion: 0.06
Nodes (27): deferredBinary, keptEntries, BinaryEntry, TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), checkBtrfsKeys(), ensureBtrfsLayout(), BinariesSummary() (+19 more)

### Community 11 - "policyEditor"
Cohesion: 0.12
Nodes (12): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), netWarnings(), newNetGuardTab() (+4 more)

### Community 12 - "protectedPage"
Cohesion: 0.06
Nodes (20): backendLink, checkGrid, configView, forwardForm, forwardView, protectedPage, recordReview, recordView (+12 more)

### Community 13 - "plan.go"
Cohesion: 0.35
Nodes (4): recordTick, isAbsoluteCandidate(), TestCatalogSanity(), hint

### Community 14 - "GuardInodeKey"
Cohesion: 0.07
Nodes (28): binaryVerifyState, refusedReplacements, appendInodeRow(), batchInodeRows(), inodeRowsOf(), newVerifyState(), supersedeUnnamed(), appletTag() (+20 more)

### Community 15 - "editView"
Cohesion: 0.12
Nodes (8): codeEntry, editView, Entry, selectable(), fileSize(), joinRel(), newCodeEntry(), relOf()

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (9): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent (+1 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (57): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+49 more)

### Community 18 - "wipe"
Cohesion: 0.18
Nodes (20): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), openFileVaultWithMasterKey(), recoverFileInPlace(), rejectSidecarSymlink() (+12 more)

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (12): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, dirEntry (+4 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.14
Nodes (49): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+41 more)

### Community 22 - "Resource"
Cohesion: 0.06
Nodes (11): fakeDaemon, findResource(), Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, uniqueEncryptionRoots() (+3 more)

### Community 23 - "auth_test.go"
Cohesion: 0.13
Nodes (40): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+32 more)

### Community 24 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (72): inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), current_is_inspector() (+64 more)

### Community 26 - "TrustGuard"
Cohesion: 0.07
Nodes (18): tempGrant, tempMaskOp, tempRow, trustRows, deleteInoKeys(), deleteResKeys(), setSupersedeTrust(), syncU32Map() (+10 more)

### Community 27 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.10
Nodes (26): class, pos, span, classOf(), column(), dedent(), Edit(), leadingIndent() (+18 more)

### Community 28 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.04
Nodes (45): recordTickCmd(), recordModel, diffModel, editorModel, Stats, parseKBValue(), parseStat(), parseStatus() (+37 more)

### Community 29 - "shellrc.go"
Cohesion: 0.16
Nodes (35): setupBunTmpdir(), rcTarget, BunTmpDir(), bunTargets(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv(), EnsureBunTmpDir(), ensureRCBlock() (+27 more)

### Community 30 - "Config"
Cohesion: 0.07
Nodes (59): askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), appendSectionsAndEdit(), applyDiffAdditions(), collectDiffAdditions(), configPaths() (+51 more)

### Community 31 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.18
Nodes (12): NewGuardUseCase(), NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle() (+4 more)

### Community 32 - "UI"
Cohesion: 0.09
Nodes (19): showChangelog(), confirmUpdate(), isTerminal(), TestConfirmUpdateRoutesToActiveUI(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), cancelUI, TestEditTextMapsUICancel(), MultiSelectKeymap() (+11 more)

### Community 33 - "NetEventType"
Cohesion: 0.06
Nodes (19): snapNetMonitor, NetBpfEvent, NetEvent, NetEventType, eventsOrAll(), newNetEventChecks(), Cstr(), FormatAddr() (+11 more)

### Community 34 - "picker"
Cohesion: 0.13
Nodes (7): keyEntry, picker, pickMode, errText(), isFindShortcut(), newKeyEntry(), newPickRow()

### Community 35 - "fileedit_model_test.go"
Cohesion: 0.18
Nodes (27): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+19 more)

### Community 37 - "backend"
Cohesion: 0.11
Nodes (8): checkSecret(), File, backend, backendOp, backendState, Message, Request, RequireDaemonStopped()

### Community 38 - "controlServer"
Cohesion: 0.10
Nodes (22): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readBareRequest(), readClientLines(), readForwardRequest(), readGrantRequest(), readPathLines() (+14 more)

### Community 39 - "uninstall/uninstall.go"
Cohesion: 0.05
Nodes (54): confirmMasterKeyOverwrite(), runGenKey(), TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), askDecrypt(), cleanOrphanedMetadata(), decryptDirectories() (+46 more)

### Community 40 - "ConfinedEntry"
Cohesion: 0.12
Nodes (29): HashStore, sharedHashEntry, ConfinedEntry(), claimHash(), hashFileShared(), hashPathShared(), HashStats(), publishHash() (+21 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 42 - "absPath"
Cohesion: 0.07
Nodes (7): IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath(), TestMain(), helperChild(), TestMain()

### Community 44 - "desktop.go"
Cohesion: 0.07
Nodes (47): installHooks(), releaseBuild(), Run(), staged(), installGUI(), needsBinary(), RunWith(), TestNeedsBinary() (+39 more)

### Community 45 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 46 - "daemon_multicall_test.go"
Cohesion: 0.17
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 48 - "safeio.go"
Cohesion: 0.06
Nodes (42): applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig(), putConfig(), runEditConfig(), TestDaemonUnitWithMetadataOutput() (+34 more)

### Community 49 - "fadeScroll"
Cohesion: 0.08
Nodes (16): fadeBar, fadeBarRenderer, fadeScroll, fadeScrollRenderer, fadeScrollOf(), newFadeScroll(), newTestFadeScroll(), settle() (+8 more)

### Community 50 - "install_ui.go"
Cohesion: 0.11
Nodes (11): installUI, multiAnswer, resizeWatch, sizedDialog, textAnswer, scrollable(), ask(), captureLogs() (+3 more)

### Community 51 - "cursor.go"
Cohesion: 0.12
Nodes (14): buttonCursor, handCursors(), handCursorsOn(), newButtonCursor(), showConfirm(), showCustomConfirm(), showDialog(), showErrorDialog() (+6 more)

### Community 52 - "Client"
Cohesion: 0.09
Nodes (14): Client, childEnv(), launch(), Start(), TestChildEnvIsAllowlisted(), TestHelperBackend(), TestReaderCapsLines(), Reader (+6 more)

### Community 53 - "watchGroup"
Cohesion: 0.11
Nodes (36): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+28 more)

### Community 54 - "tabs.go"
Cohesion: 0.03
Nodes (22): appTheme, clickPanel, iconTileRenderer, navItem, navItemRenderer, panelRenderer, pillRenderer, tabBar (+14 more)

### Community 55 - ".apply"
Cohesion: 0.12
Nodes (15): Plan, sectionMaterialized(), TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(), TestApplyRollsBackOnRefusal(), testPlan(), TestPlanAdmitsOnlyWhatWasRecorded() (+7 more)

### Community 56 - "prepareDaemonStart"
Cohesion: 0.08
Nodes (25): CheckBPFLSM(), CheckEBPF(), ParseNetEventsFlag(), configureDaemonLogging(), newPinGeneration(), prepareDaemonStart(), resolveConfigPath(), runBPFCheck() (+17 more)

### Community 57 - "mutedText"
Cohesion: 0.09
Nodes (31): filesPage, recordForm, newLogo(), newBinaryList(), newBinaryRow(), newFadeVScroll(), newTargetRow(), newCounterStrip() (+23 more)

### Community 58 - "engine"
Cohesion: 0.05
Nodes (29): engine, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), TestSetInspectorsBounded() (+21 more)

### Community 59 - "os.File"
Cohesion: 0.04
Nodes (68): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), ownerPaths(), checkFreshDir(), openDirNoFollow(), renewDir(), setConfinedHomes() (+60 more)

### Community 60 - "CandidateDir"
Cohesion: 0.07
Nodes (28): entryWriters(), reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName() (+20 more)

### Community 61 - "Recorder"
Cohesion: 0.16
Nodes (15): drain(), Recorder, Target, monitorTargets(), newRecorder(), selfExe(), sortedKeys(), Start() (+7 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.14
Nodes (20): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+12 more)

### Community 63 - "Network"
Cohesion: 0.10
Nodes (10): binaryList, netGuardPill(), guardRow(), clonePolicy(), NetExtra, NetPolicy, Network, verifyHashes() (+2 more)

### Community 64 - "Binary"
Cohesion: 0.11
Nodes (28): lines(), pickBinaries(), pickLibs(), reviewRecording(), appletLine(), Classify(), classifyGroup(), classifyLibs() (+20 more)

### Community 65 - "pinstate.go"
Cohesion: 0.27
Nodes (15): pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID(), TestPinStateRoundTrip(), TestReadPinStateMalformed() (+7 more)

### Community 66 - "wm.go"
Cohesion: 0.10
Nodes (27): invokingIDs(), invokingUID(), detectWM(), reachable(), centered(), Detect(), FitSize(), Session (+19 more)

### Community 67 - ".FilterExistingWhitelist"
Cohesion: 0.11
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 68 - "recordFlow"
Cohesion: 0.14
Nodes (15): encryptBlocker(), newRecordFlow(), confirmAndApply(), offerBackupDeletion(), pickEncryption(), runRecord(), runRecordView(), recordFlow (+7 more)

### Community 69 - "lockAndDeprovision"
Cohesion: 0.14
Nodes (9): deprovisionKind, migrationTarget, classifyDeprovision(), classifyDeprovisionErr(), applyRawKeyPolicy(), classifyMigrationTarget(), isDeprovisionBusy(), isDeprovisionMissing() (+1 more)

### Community 70 - "fcntl"
Cohesion: 0.10
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 72 - "Guard"
Cohesion: 0.08
Nodes (11): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), exeRow, rootHandle, openRoot(), SharedPinDegraded() (+3 more)

### Community 73 - "FileTarget"
Cohesion: 0.14
Nodes (10): Target, guardPill(), Files, FileSpec, FileTarget, Policy, newFileTarget(), NormalizeSpec() (+2 more)

### Community 74 - "KernelDev"
Cohesion: 0.15
Nodes (18): fdStat, physicalParent(), statFD(), KernelDev(), LstatInode(), BtrfsDevices(), btrfsSysfsDevices(), TestBtrfsSysfsDevices() (+10 more)

### Community 76 - "catalogRefresher"
Cohesion: 0.15
Nodes (9): drain(), newCatalogRefresher(), readEvents(), swapFile(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), startCatalogRefresh(), catalogRefresher (+1 more)

### Community 77 - "update_test.go"
Cohesion: 0.11
Nodes (14): fetchReleases(), rsaPublicKeyPEM(), signedChecksum(), TestFetchReleases(), TestFetchReleasesHTTPError(), TestFilterChannel(), testKeyPair(), TestNewerThanInstalled() (+6 more)

### Community 78 - "compileGlobNames"
Cohesion: 0.17
Nodes (9): GlobChild, statFunc, compileGlobNames(), globText(), resolveBits(), resolveChildren(), statKey(), TestCompileGlobNames_RepeatedPatternOrsBits() (+1 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.40
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "pickfs.go"
Cohesion: 0.27
Nodes (13): listOpts, pickEntry, treeSearch, TestSearchTreeFollowsNoSymlink(), plural(), entryIsDir(), filterEntries(), listDir() (+5 more)

### Community 81 - "newTestTarget"
Cohesion: 0.09
Nodes (23): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+15 more)

### Community 82 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 83 - "changelog.go"
Cohesion: 0.27
Nodes (16): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), compareStableVersions(), filterChannel(), findUpdate(), latestTime() (+8 more)

### Community 84 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 86 - "render.go"
Cohesion: 0.12
Nodes (13): lineStyles, viewportState, baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces(), setGutter() (+5 more)

### Community 87 - "uninstall_system.go"
Cohesion: 0.17
Nodes (18): confirmKeyDeletion(), run(), runUninstall(), RunWith(), detectSSHAgentUnits(), isInstallerSSHAgentUnit(), removeKeyAndEmptyDir(), removeMasterKey() (+10 more)

### Community 88 - "liveSession"
Cohesion: 0.17
Nodes (13): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), forwardCheck() (+5 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "lockRootRecovering"
Cohesion: 0.16
Nodes (15): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), ensurePinStateFilePlaceholder(), ensurePlaceholder(), bootstrapSelfProtectPlaceholders(), readTempJournal(), stripTempJournal() (+7 more)

### Community 92 - "Monitor"
Cohesion: 0.16
Nodes (5): NewMonitor(), putBtrfsLayout(), RealFileEngines(), Monitor, watchedInode

### Community 93 - "netinject.c"
Cohesion: 0.12
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.08
Nodes (14): trustHook, foreignPinnedLink(), reasonText(), eventTypeKey(), mcNameKey(), statInodeKey(), judgePasses(), parseProcs() (+6 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "io.Writer"
Cohesion: 0.19
Nodes (12): formatElapsed(), printBinary(), printRecording(), candidates(), confirmable(), confirmAll(), describe(), hashNow() (+4 more)

### Community 100 - "highlight_test.go"
Cohesion: 0.13
Nodes (26): New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep(), TestCursorDoesNotBlink(), TestEditKeys(), TestEnterKeepsIndentation(), TestEnterPastNinetyNineLines(), TestSetTextStartsAtFirstLine() (+18 more)

### Community 101 - "downloadFile"
Cohesion: 0.24
Nodes (7): downloadFile(), downloadReleaseFiles(), fetchRelease(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleaseByTag(), progressReader

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "installPage"
Cohesion: 0.04
Nodes (38): actionGate, binaryRow, binSource, handButton, iconTile, installAction, installedBinary, InstallMode (+30 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.22
Nodes (13): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec() (+5 more)

### Community 106 - "recordops.c"
Cohesion: 0.45
Nodes (13): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+5 more)

### Community 107 - "tempGrantSetup"
Cohesion: 0.29
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.17
Nodes (9): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), btrfsTypes(), TestBtrfsLayoutFromNestedMembers() (+1 more)

### Community 110 - "hashCache"
Cohesion: 0.12
Nodes (17): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, BinaryStat, DurableStat, HashCache (+9 more)

### Community 111 - "mountTable"
Cohesion: 0.18
Nodes (9): MountInfo, mountinfoUnitTest, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo() (+1 more)

### Community 112 - "NetGuardEvent"
Cohesion: 0.07
Nodes (26): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runHeadless(), runNetworkGuard(), runTUI(), runHeadless() (+18 more)

### Community 113 - "buildTrustedSet"
Cohesion: 0.21
Nodes (13): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+5 more)

### Community 115 - ".step"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 116 - "EventType"
Cohesion: 0.06
Nodes (33): EventType, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), binaryStatOf(), TempBinary, TemporaryGrant, TempRule, loadPinnedExeMaps() (+25 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - ".TestRecord_GuardParity"
Cohesion: 0.36
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 121 - "GenerateConf"
Cohesion: 0.05
Nodes (59): logRefreshChange(), TestLogRefreshChange(), LibraryBlock, RefreshOptions, SectionScan, findSectionEnd(), findSectionStart(), firstContentLine() (+51 more)

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

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - ".launch"
Cohesion: 0.31
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "install/install.go"
Cohesion: 0.06
Nodes (71): applyLiveRefresh(), buildBinaryIfNeeded(), certifyInstall(), checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), encryptDirectories(), ensureInstalledBinary(), ensureMasterKey() (+63 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "dialLiveSession"
Cohesion: 0.09
Nodes (23): ParseEventsFlag(), reserveStdout(), runBackend(), dialLiveSession(), LiveModeAvailable(), TestValidateFlagsEditConfig(), runEditProtected(), validateFlags() (+15 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "GuardEvent"
Cohesion: 0.04
Nodes (25): holdHeadless(), HeadlessLine(), runGuardHeadless(), bpfGuardEvent, snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestLocalizePaths() (+17 more)

### Community 143 - "time.Duration"
Cohesion: 0.17
Nodes (12): procStats, statsSample, cpuTime(), humanBytes(), rssBytes(), New(), recv(), TestSubDeliversMissedTickLate() (+4 more)

### Community 144 - "auditGroupCoverage"
Cohesion: 0.16
Nodes (14): auditAfterEdit(), auditEntry(), auditFindings(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource() (+6 more)

### Community 145 - "newTestVetter"
Cohesion: 0.29
Nodes (11): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+3 more)

### Community 146 - "session_test.go"
Cohesion: 0.19
Nodes (13): ancestorEnv(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField(), TestDisplayValue(), TestOwnedByChecksTypeAndOwner(), TestReadOwnedFileRefusesSymlink() (+5 more)

### Community 147 - "sync.Mutex"
Cohesion: 0.07
Nodes (13): logGrid, logHook, logLine, runView, stageItem, stageState, logRow(), logStyle() (+5 more)

### Community 148 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 149 - "fscrypt/orphans.go"
Cohesion: 0.19
Nodes (20): modifiedContextWithSource(), appProtectorSet(), CleanOrphanedMetadata(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), LivePolicyDescriptors() (+12 more)

### Community 150 - "watchPattern"
Cohesion: 0.27
Nodes (10): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), addSystemBinary(), catalogPatterns(), catalogWatchPlan(), inHome(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors() (+2 more)

### Community 152 - "New"
Cohesion: 0.12
Nodes (28): TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), askEncryption(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic() (+20 more)

### Community 154 - "buildOneGuard"
Cohesion: 0.06
Nodes (51): backingDeviceUnion(), buildConcurrency(), buildGuards(), buildOneGuard(), eventFilterOptions(), reloadOnce(), reloadSlotsNeeded(), startGuardedDaemon() (+43 more)

### Community 155 - "layout.go"
Cohesion: 0.43
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 157 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "backend_test.go"
Cohesion: 0.27
Nodes (19): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+11 more)

### Community 161 - "vaultfs.go"
Cohesion: 0.09
Nodes (26): entryOf(), DefaultNewFileMode(), FileModeToUnixMode(), Vault, IsBinaryContent(), IsOctalMode(), Open(), ParseMode() (+18 more)

### Community 162 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 163 - "runDaemon"
Cohesion: 0.09
Nodes (19): newControlManager(), startControlManager(), catchLifecycleSignals(), makeReloadHandler(), markReady(), noteReloadedSettings(), notifySystemdReady(), runDaemon() (+11 more)

### Community 164 - "FileEvent"
Cohesion: 0.11
Nodes (7): snapMonitor, BpfEvent, eventUnitTest, FileEvent, DecodeBpfEvent(), fakeSource, fakeMonitorRepo

### Community 165 - "feed.go"
Cohesion: 0.06
Nodes (44): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+36 more)

### Community 166 - "shell"
Cohesion: 0.04
Nodes (52): fakeLink, homePage, pageDef, pageID, reel, runHost, shell, snapEngines (+44 more)

### Community 167 - "runMonitor"
Cohesion: 0.11
Nodes (18): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+10 more)

### Community 168 - "Ledger"
Cohesion: 0.11
Nodes (5): applySchema(), Ledger, Pending, Source, fakeLedger

### Community 170 - "find.go"
Cohesion: 0.11
Nodes (27): findButton, findAll(), lowerRunes(), NewFinder(), newFindInput(), overlay(), runesEqual(), findEditor() (+19 more)

### Community 171 - "DiscoverInfraBinaries"
Cohesion: 0.33
Nodes (6): DiscoverInfraBinaries(), infraFromRunning(), runningExecutables(), pathIsRunning(), TestDiscoverInfraBinaries(), TestInfraFromRunning()

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 176 - "dialog_test.go"
Cohesion: 0.45
Nodes (10): openPicker(), pickerTree(), shownNames(), TestPickerArrowHighlightsFolderEnterOpensIt(), TestPickerCtrlFFocusesSearch(), TestPickerFolderModeChoosesShownFolder(), TestPickerSearchFiltersFolder(), TestPickerSubfolderSearchPicksFile() (+2 more)

### Community 177 - "copyRegularAt"
Cohesion: 0.15
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 179 - "fileedit_symlink_test.go"
Cohesion: 0.42
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.22
Nodes (9): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener (+1 more)

### Community 182 - "github.com/spf13/cobra.Command"
Cohesion: 0.07
Nodes (38): AddServeFlags(), credentialFlagState(), ServeConfig, isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+30 more)

### Community 186 - "Binary trust"
Cohesion: 0.29
Nodes (7): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

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

### Community 195 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 196 - "edit-protected"
Cohesion: 0.22
Nodes (9): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Recording a new path: `--record`, Temporary access for other programs: `--forward` (+1 more)

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - ".refreshOnce"
Cohesion: 0.43
Nodes (6): sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad()

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 202 - "DaemonEvent"
Cohesion: 0.09
Nodes (36): UIDResolver, NewUIDResolver(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI(), runTUI() (+28 more)

### Community 203 - "TestBatchInodeRows"
Cohesion: 0.40
Nodes (6): inodeBatchFn, fakeBatches(), TestBatchInodeRows(), countInodeRows(), InodeOccupancy(), TestCountInodeRows()

### Community 204 - "vouchedDevs"
Cohesion: 0.33
Nodes (6): isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed(), vouchedDevs()

### Community 205 - "parsePasswd"
Cohesion: 0.50
Nodes (4): parsePasswd(), shell(), TestParsePasswd(), TestParsePasswdMalformed()

### Community 206 - "Desktop GUI"
Cohesion: 0.15
Nodes (13): Actions, Desktop GUI, File and folder pickers, Files, Home, Install, Network, Protected (+5 more)

### Community 209 - "WithBottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 216 - "LibraryClosure"
Cohesion: 0.07
Nodes (30): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), isScript(), ldConfInclude(), ldSoConfDirs() (+22 more)

### Community 220 - "update/update.go"
Cohesion: 0.09
Nodes (43): applyUpdate(), archAssets(), assetsFor(), channelOf(), checkUpdatePreconditions(), confirmUnit(), deploy(), deployGUI() (+35 more)

### Community 221 - "Daemon configuration"
Cohesion: 0.20
Nodes (9): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Process inspectors, Reloading, Several guarded trees, one vault (+1 more)

### Community 224 - "time.Time"
Cohesion: 0.13
Nodes (11): LockedError, isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet, sweepID (+3 more)

### Community 225 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 230 - "headless.go"
Cohesion: 0.36
Nodes (6): Check(), CheckGUI(), inspect(), buildTiny(), TestCheck(), TestCheckGUI()

### Community 237 - "Serve"
Cohesion: 0.08
Nodes (29): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+21 more)

## Knowledge Gaps
- **31 isolated node(s):** `IntegrationSuite`, `github.com/Virgula0/app-listener`, `open_by_handle_at POC`, `statonly stat/statx metadata leak`, `app-listener-gui best-effort build` (+26 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 586 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **38 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `daemon/daemon.go`, `.launch`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `daemon_multicall_test.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `.TestRecord_GuardParity`, `IntegrationSuite`?**
  _High betweenness centrality (0.042) - this node is a cross-community bridge._
- **What connects `IntegrationSuite`, `github.com/Virgula0/app-listener`, `open_by_handle_at POC` to the rest of the system?**
  _31 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.10996879954989515 - nodes in this community are weakly interconnected._
- **Why does `Editor keys` connect `edit-protected` to `desktop.go`, `backend`?**
  _High betweenness centrality (0.041) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.022289766970618033 - nodes in this community are weakly interconnected._
- **Why does `edit-protected` connect `edit-protected` to `README.md`?**
  _High betweenness centrality (0.041) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.08958130477117819 - nodes in this community are weakly interconnected._